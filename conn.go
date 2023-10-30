package zysms

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/enum"
	"github.com/zhiyin2021/zysms/smserror"
	"github.com/zhiyin2021/zysms/utils"
)

type sms_conn struct {
	Data any
	// Logger *logrus.Entry
	UUID           string
	ctx            context.Context
	stop           func()
	activeCount    int32
	activeInterval int
	IsHealth       bool
	nodeId         uint32

	net.Conn
	State          enum.State
	Protocol       codec.SmsProto
	Typ            codec.Version
	counter        int32
	logger         *logrus.Entry
	activeFail     int32
	extParam       map[string]string
	checkVer       bool
	autoActiveResp bool

	action sms_action
	mutex  sync.Mutex
}

type sms_action interface {
	login(uid, pwd string) error
	logout()
	recv() (codec.PDU, error)
	active_test() error
}

func newConn(conn net.Conn, proto codec.SmsProto) *sms_conn {
	c := &sms_conn{
		Conn:           conn,
		UUID:           utils.RandomStr(10),
		Typ:            proto.Version(),
		Protocol:       proto,
		logger:         logrus.WithFields(logrus.Fields{"r": conn.RemoteAddr(), "v": proto.String()}),
		extParam:       map[string]string{},
		checkVer:       false,
		autoActiveResp: true,
		activeCount:    0,
		activeInterval: 5,
	}

	c.ctx, c.stop = context.WithCancel(context.Background())
	c.SetState(enum.CONN_CONNECTED)

	switch proto {
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
	return c
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

func (c *sms_conn) startActiveTest(errEvent func(Conn, error), heartbeatNoResp func(Conn, int)) {
	c.IsHealth = true
	if c.activeInterval > 0 {
		tryGO(func() {
			t := time.NewTicker(time.Duration(c.activeInterval) * time.Second)
			defer t.Stop()
			for {
				select {
				case <-c.ctx.Done():
					// once conn close, the goroutine should exit
					return
				case <-t.C:
					n, err := c.sendActiveTest()
					if err != nil {
						errEvent(c, fmt.Errorf("心跳请求异常:%s", err))
						return
					}
					if c.activeCount > 0 && n >= c.activeCount {
						if errEvent != nil {
							errEvent(c, fmt.Errorf("间隔(%ds),%d次心跳异常,关闭连接", c.activeCount, c.activeInterval))
						} else {
							c.Logger().Errorf("间隔(%ds),%d次心跳异常,关闭连接", c.activeCount, c.activeInterval)
						}
						c.Close()
						return
					} else if heartbeatNoResp != nil {
						heartbeatNoResp(c, int(n))
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
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.stop()
	if c.State != enum.CONN_CLOSED {

		if c.State == enum.CONN_AUTHOK {
			if c.action != nil {
				c.action.logout()
			}
			time.Sleep(100 * time.Millisecond)
		}
		c.Conn.Close() // close the underlying net.Conn
		c.State = enum.CONN_CLOSED
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
			c.activeInterval = utils.MapItem(ext, "active_interval", int(5))

			c.checkVer = utils.MapItem(ext, "check_version", 0) == 1
			c.autoActiveResp = utils.MapItem(ext, "auto_active_resp", 1) == 1
		})
	}
}

func (c *sms_conn) SetState(state enum.State) {
	c.State = state
}

// SendPkt pack the smpp packet structure and send it to the other peer.
func (c *sms_conn) SendPDU(pdu PDU) error {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorln("smpp.send.panic:", err)
			c.Close()
		}
	}()
	if c.State == enum.CONN_CLOSED {
		c.Close()
		return smserror.ErrConnIsClosed
	}
	if pdu == nil {
		return smserror.ErrPktIsNil
	}
	buf := codec.NewWriter()
	c.Logger().Debugf("send pdu:%T , %d , %d", pdu, c.Typ, buf.Len())
	pdu.Marshal(buf)
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
