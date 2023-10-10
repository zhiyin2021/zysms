package zysms

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms/cmpp"
	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/enum"
	"github.com/zhiyin2021/zysms/smserror"
	"github.com/zhiyin2021/zysms/utils"
)

type cmppConn struct {
	net.Conn
	State      enum.State
	Typ        codec.Version
	counter    int32
	logger     *logrus.Entry
	checkVer   bool
	activeFail int32
	extParam   map[string]string
	// activePeer bool // 默认false,当前连接发送心跳请求, 当收到对方心跳请求后,设置true,不再发送心跳请求
	// activeLast time.Time
}

// New returns an abstract structure for successfully
// established underlying net.Conn.
func newCmppConn(conn net.Conn, typ codec.Version) smsConn {
	c := &cmppConn{
		Conn:     conn,
		Typ:      typ,
		logger:   logrus.WithFields(logrus.Fields{"r": conn.RemoteAddr()}),
		extParam: map[string]string{},
		checkVer: false,
	}
	return c
}
func (c *cmppConn) Ver() codec.Version {
	return c.Typ
}
func (c *cmppConn) Auth(uid string, pwd string) error {
	// Login to the server.
	req := cmpp.NewConnReq(c.Typ).(*cmpp.ConnReq)
	req.SrcAddr = uid
	req.Secret = pwd
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

	if rsp, ok := p.(*cmpp.ConnResp); ok {
		if c.checkVer && rsp.Version != c.Typ {
			return smserror.ErrVersionNotMatch
		}
		status = uint8(rsp.Status)
	} else {
		return smserror.ErrRespNotMatch
	}

	if status != 0 {
		return smserror.NewSmsErr(int(status), "cmpp.login.error")
	}
	c.SetState(enum.CONN_AUTHOK)
	return nil
}
func (c *cmppConn) close() {
	if c != nil {
		if c.State == enum.CONN_CLOSED {
			return
		}
		c.Conn.Close() // close the underlying net.Conn
		c.State = enum.CONN_CLOSED
	}
}

func (c *cmppConn) SetState(state enum.State) {
	c.State = state
}

func (c *cmppConn) setExtParam(ext map[string]string) {
	if ext != nil {
		c.checkVer = utils.MapItem(ext, "check_version", 0) == 1
		c.extParam = ext
	}
}

// SendPkt pack the smpp packet structure and send it to the other peer.
func (c *cmppConn) SendPDU(pdu codec.PDU) error {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorln("cmpp.send.panic:", err)
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

func (c *cmppConn) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
func (c *cmppConn) Logger() *logrus.Entry {
	return c.logger
}
func (c *cmppConn) SetReadDeadline(timeout time.Duration) {
	if timeout > 0 {
		c.Conn.SetReadDeadline(time.Now().Add(timeout))
	}
}

// RecvAndUnpackPkt receives cmpp byte stream, and unpack it to some cmpp packet structure.
func (c *cmppConn) RecvPDU() (codec.PDU, error) {
	if c.State == enum.CONN_CLOSED {
		return nil, smserror.ErrConnIsClosed
	}

	pdu, err := cmpp.Parse(c.Conn, c.Typ)
	if err != nil {
		return nil, err
	}

	switch p := pdu.(type) {
	case *cmpp.ActiveTestReq: // 当收到心跳请求,内部直接回复心跳,并递归继续获取数据
		resp := p.GetResponse()
		c.SendPDU(resp)
		// if !c.activePeer {
		// 	c.activePeer = true
		// }
	case *cmpp.ActiveTestResp: // 当收到心跳回复,内部直接处理,并递归继续获取数据
		atomic.AddInt32(&c.counter, -1)
		// c.activeLast = time.Now()
	case *cmpp.ConnResp: // 当收到登录回复,内部先校验版本
		if c.checkVer && p.Version != c.Typ {
			return nil, fmt.Errorf("cmpp version not match [ local: %d != remote: %d ]", c.Typ, p.Version)
		}
	case *cmpp.ConnReq: // 当收到登录回复,内部先校验版本
		// 服务端自适应版本
		c.Typ = p.Version
		c.logger = logrus.WithFields(logrus.Fields{"r": c.RemoteAddr(), "v": c.Typ})
	}
	return pdu, nil

}
func (c *cmppConn) sendActiveTest() (int32, error) {
	p := cmpp.NewActiveTestReq(c.Typ)
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

type cmppListener struct {
	net.Listener
	extParam map[string]string
}

func newCmppListener(l net.Listener, extParam map[string]string) *cmppListener {
	return &cmppListener{l, extParam}
}

func (l *cmppListener) accept() (smsConn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	tc := c.(*net.TCPConn)
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(30 * time.Second) // 1min

	conn := newCmppConn(c, cmpp.V30)
	conn.SetState(enum.CONN_CONNECTED)
	return conn, nil
}

func (l *cmppListener) Close() error {
	return l.Listener.Close()
}
