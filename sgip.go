package zysms

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/enum"
	"github.com/zhiyin2021/zysms/sgip"
	"github.com/zhiyin2021/zysms/smserror"
	"github.com/zhiyin2021/zysms/utils"
)

type sgipConn struct {
	net.Conn
	State      enum.State
	Typ        codec.Version
	counter    int32
	logger     *logrus.Entry
	checkVer   bool
	nodeId     uint32
	activeFail int32
	extParam   map[string]string
	// activePeer bool // 默认false,当前连接发送心跳请求, 当收到对方心跳请求后,设置true,不再发送心跳请求
}

// New returns an abstract structure for successfully
// established underlying net.Conn.
func newSgipConn(conn net.Conn, typ codec.Version) smsConn {
	c := &sgipConn{
		Conn:     conn,
		Typ:      typ,
		logger:   logrus.WithFields(logrus.Fields{"r": conn.RemoteAddr()}),
		extParam: map[string]string{},
		checkVer: false,
	}
	return c
}
func (c *sgipConn) Ver() codec.Version {
	return c.Typ
}
func (c *sgipConn) Auth(uid string, pwd string) error {
	// Login to the server.
	req := sgip.NewBindReq(c.Typ, c.nodeId).(*sgip.BindReq)
	req.LoginName = uid
	req.LoginPassword = pwd
	req.LoginType = 1
	req.Version = c.Typ

	err := c.SendPDU(req)
	if err != nil {
		return err
	}
	p, err := c.RecvPDU()
	if err != nil {
		return err
	}
	var status uint8

	if rsp, ok := p.(*sgip.BindResp); ok {
		if c.checkVer && rsp.Version != c.Typ {
			return smserror.ErrVersionNotMatch
		}
		status = uint8(rsp.Status)
	} else {
		return smserror.ErrRespNotMatch
	}

	if status != 0 {
		// if status <= smserror.ErrnoConnOthers { //ErrnoConnOthers = 5
		// 	err = smserror.ConnRspStatusErrMap[status]
		// } else {
		// 	err = smserror.ConnRspStatusErrMap[smserror.ErrnoConnOthers]
		// }
		return smserror.NewSmsErr(int(status), "sgip.login.error")
	}
	c.SetState(enum.CONN_AUTHOK)
	return nil
}
func (c *sgipConn) close() {
	if c != nil {
		if c.State == enum.CONN_CLOSED {
			return
		}
		if c.State == enum.CONN_AUTHOK {
			c.SendPDU(sgip.NewUnbindReq(c.Typ, c.nodeId))
			time.Sleep(100 * time.Millisecond)
		}
		c.Conn.Close() // close the underlying net.Conn
		c.State = enum.CONN_CLOSED
	}
}

func (c *sgipConn) SetState(state enum.State) {
	c.State = state
}

func (c *sgipConn) setExtParam(ext map[string]string) {
	if ext != nil {
		c.checkVer = utils.MapItem(ext, "check_version", 0) == 1
		c.nodeId = utils.MapItem(ext, "node_id", uint32(0))
		c.extParam = ext
	}
}

// SendPkt pack the smpp packet structure and send it to the other peer.
func (c *sgipConn) SendPDU(pdu codec.PDU) error {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorln("sgip.send.panic:", err)
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

func (c *sgipConn) Logger() *logrus.Entry {
	return c.logger
}
func (c *sgipConn) SetReadDeadline(timeout time.Duration) {
	if timeout > 0 {
		c.Conn.SetReadDeadline(time.Now().Add(timeout))
	}
}

// RecvAndUnpackPkt receives sgip byte stream, and unpack it to some sgip packet structure.
func (c *sgipConn) RecvPDU() (codec.PDU, error) {
	if c.State == enum.CONN_CLOSED {
		return nil, smserror.ErrConnIsClosed
	}

	pdu, err := sgip.Parse(c.Conn, c.Typ, c.nodeId)
	if err != nil {
		return nil, err
	}

	switch p := pdu.(type) {

	case *sgip.ReportReq: // 当收到心跳请求,内部直接回复心跳,并递归继续获取数据
		resp := p.GetResponse().(*sgip.ReportResp)
		resp.Status = 0
		c.SendPDU(resp)
	case *sgip.ReportResp: // 当收到心跳回复,内部直接处理,并递归继续获取数据
		atomic.AddInt32(&c.counter, -1)
	case *sgip.BindResp: // 当收到登录回复,内部先校验版本
		if c.checkVer && p.Version != c.Typ {
			return nil, fmt.Errorf("sgip version not match [ local: %d != remote: %d ]", c.Typ, p.Version)
		}
	case *sgip.BindReq: // 当收到登录回复,内部先校验版本
		// 服务端自适应版本
		c.Typ = p.Version
		c.logger = logrus.WithFields(logrus.Fields{"r": c.RemoteAddr(), "v": c.Typ})
	}
	return pdu, nil

}

func (c *sgipConn) sendActiveTest() (int32, error) {
	p := sgip.NewReportReq(c.Typ, c.nodeId).(*sgip.ReportReq)
	err := c.SendPDU(p)
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

type sgipListener struct {
	net.Listener
	extParam map[string]string
}

func newSgipListener(l net.Listener, extParam map[string]string) *sgipListener {
	return &sgipListener{l, extParam}
}

func (l *sgipListener) accept() (smsConn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	tc := c.(*net.TCPConn)
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(30 * time.Second) // 1min

	conn := newSgipConn(c, sgip.V12)
	conn.SetState(enum.CONN_CONNECTED)
	return conn, nil
}

func (l *sgipListener) Close() error {
	return l.Listener.Close()
}
