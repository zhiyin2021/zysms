package zysms

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zycli/cache"
	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/enum"
	"github.com/zhiyin2021/zysms/smserror"
	"github.com/zhiyin2021/zysms/utils"
)

type activeTestItem struct {
	time  time.Time
	timer *time.Timer
	flag  int32
}
type sms_conn struct {
	Data any
	// Logger *logrus.Entry
	sid            string
	ctx            context.Context
	stop           func()
	activeCount    int32
	activeInterval int32
	IsHealth       bool
	nodeId         uint32

	net.Conn
	Protocol       codec.SmsProto
	Typ            codec.Version
	counter        int32
	logger         *logrus.Entry
	activeFail     int32
	extParam       map[string]string
	checkVer       bool
	autoActiveResp bool

	action sms_action
	delay  *utils.Queue

	Connected int32
	IsAuth    bool
	cache     *cache.Memory

	errEvent        func(Conn, error)
	heartbeatNoResp func(Conn, int)

	parent *SMS
}

type sms_action interface {
	login(uid, pwd string) error
	logout()
	recv() (codec.PDU, error)
	active_test() error
}

func newConn(conn net.Conn, parent *SMS) *sms_conn {
	sid := utils.Md5(fmt.Sprintf("%s%s%d", conn.RemoteAddr(), conn.LocalAddr(), time.Now().UnixNano()))[8:24]
	addr := fmt.Sprintf("%s->%s", conn.LocalAddr(), conn.RemoteAddr())
	c := &sms_conn{
		Conn:           conn,
		sid:            sid,
		Typ:            parent.proto.Version(),
		Protocol:       parent.proto,
		logger:         logrus.WithFields(logrus.Fields{"sid": sid, "addr": addr, "v": parent.proto.String()}),
		extParam:       map[string]string{},
		checkVer:       false,
		autoActiveResp: true,
		activeCount:    0,
		activeInterval: 5,
		delay:          utils.NewQueue(10),
		cache:          cache.NewMemory(time.Second * 1),
		parent:         parent,
	}
	switch parent.proto {
	case codec.CMPP20, codec.CMPP21, codec.CMPP30:
		c.action = newCmpp(c)
	case codec.SMGP30:
		c.action = newSmgp(c)
	case codec.SGIP:
		c.action = newSgip(c)
	case codec.SMPP33, codec.SMPP34:
		c.action = newSmpp(c)
	default:
		return nil
	}
	c.ctx, c.stop = context.WithCancel(context.Background())
	atomic.StoreInt32(&c.Connected, 1)
	return c
}
func (c *sms_conn) IsConnected() bool {
	return c.Connected == 1
}
func (c *sms_conn) Auth(uid string, pwd string) error {
	if c.action == nil {
		return smserror.ErrProtoNotSupport
	}
	return c.action.login(uid, pwd)
}
func (c *sms_conn) SetData(data any) {
	c.Data = data
}
func (c *sms_conn) GetData() any {
	return c.Data
}
func (c *sms_conn) SID() string {
	return c.sid
}
func (c *sms_conn) Delay() []int64 {
	return c.delay.Data()
}
func (c *sms_conn) EnabledActiveTest() {
	c.IsHealth = true
	if n := atomic.LoadInt32(&c.activeInterval); n > 0 {
		tryGO(func() {
			t := time.NewTicker(time.Duration(n) * time.Second)
			defer t.Stop()
			for {
				select {
				case <-c.ctx.Done():
					// once conn close, the goroutine should exit
					return
				case <-t.C:
					n, err := c.sendActiveTest()
					if err != nil {
						c.parent.doError(c, fmt.Errorf("心跳请求异常:%s", err))
						return
					}
					if c.activeCount > 0 && n >= c.activeCount {
						c.parent.doError(c, fmt.Errorf("间隔(%ds),%d次心跳异常,关闭连接", c.activeCount, c.activeInterval))
						c.Close()
						return
					} else if c.parent.OnHeartbeatNoResp != nil {
						c.parent.OnHeartbeatNoResp(c, int(n))
					}
				}
			}
		})
	}
}

func (c *sms_conn) sendActiveTest() (int32, error) {
	err := c.action.active_test()
	if err != nil {
		c.activeFail++
		if c.activeFail > 2 {
			return c.activeFail, err
		}
	} else {
		c.activeFail = 0
	}
	n := atomic.AddInt32(&c.counter, 1)
	return n, nil
}
func (c *sms_conn) Close() {
	if atomic.CompareAndSwapInt32(&c.Connected, enum.CONN_CONNECTED, enum.CONN_DISCONNECTED) {
		c.logger.Warnln("connection closing.")
		if c.IsAuth {
			if c.action != nil {
				c.action.logout()
			}
			time.Sleep(100 * time.Millisecond)
			c.IsAuth = false
		}
		c.stop()
		c.Conn.Close()
		c.logger.Warnln("connection closed.")
	}
}

/*
active_count 心跳失败次数
active_interval 心跳间隔
check_version 是否校验版本
system_type 系统类型[smpp 特有]
*/
func (c *sms_conn) SetExtParam(ext map[string]string) {
	if ext != nil {
		tryGO(func() {
			c.activeCount = utils.MapItem(ext, "active_count", int32(0))
			c.activeInterval = utils.MapItem(ext, "active_interval", int32(5))

			c.checkVer = utils.MapItem(ext, "check_version", 0) == 1
			c.autoActiveResp = utils.MapItem(ext, "auto_active_resp", 1) == 1
		})
	}
}

// SendPkt pack the smpp packet structure and send it to the other peer.
func (c *sms_conn) SendPDU(pdu PDU) error {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorln("smpp.send.panic:", err)
			c.Close()
		}
	}()
	if c.Connected == enum.CONN_DISCONNECTED {
		return smserror.ErrConnIsClosed
	}
	if pdu == nil {
		return smserror.ErrPktIsNil
	}
	buf := codec.NewWriter()
	pdu.Marshal(buf)
	c.Logger().Debugf("send pdu: %d , %d ,%#v ", c.Typ, buf.Len(), pdu)
	_, err := c.Conn.Write(buf.Bytes()) //block write
	if err != nil {
		c.Close()
	}
	return err
}

func (c *sms_conn) Logger() *logrus.Entry {
	return c.logger
}
func (c *sms_conn) SetReadDeadline(timeout time.Duration) {
	if timeout > 0 {
		c.Conn.SetReadDeadline(time.Now().Add(timeout))
	}
}
func (c *sms_conn) Ver() codec.Version {
	return c.Typ
}

func (c *sms_conn) activeTestReq(seq int32) {
	item := activeTestItem{
		time: time.Now(),
		flag: 0,
	}
	item.timer = time.AfterFunc(time.Second*1, func() {
		time.Sleep(100 * time.Millisecond)
		if atomic.CompareAndSwapInt32(&item.flag, 0, 1) {
			c.delay.Push(-1)
		}
	})
	c.cache.Set(fmt.Sprintf("active_test_%d", seq), &item)
}

func (c *sms_conn) activeTestResp(seq int32) {
	if tmp := c.cache.Get(fmt.Sprintf("active_test_%d", seq)); tmp != nil {
		if item, ok := tmp.(*activeTestItem); ok {
			if atomic.CompareAndSwapInt32(&item.flag, 0, 1) {
				item.timer.Stop()
				c.delay.Push(time.Since(item.time).Microseconds())
			}
		}
	}
}
