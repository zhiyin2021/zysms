package zysms

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/enum"
	"github.com/zhiyin2021/zysms/smgp"
	"github.com/zhiyin2021/zysms/smserror"
)

type smgpConn struct {
	net.Conn
	State enum.State
	Typ   codec.Version
	// for SeqId generator goroutine
	// SeqId  <-chan uint32
	// done   chan<- struct{}
	stop     func()
	ctx      context.Context
	counter  int32
	logger   *logrus.Entry
	checkVer bool
}

// New returns an abstract structure for successfully
// established underlying net.Conn.
func newSmgpConn(conn net.Conn, typ codec.Version, checkVer bool) *Conn {
	c := &smgpConn{
		Conn:     conn,
		Typ:      typ,
		logger:   logrus.WithFields(logrus.Fields{"r": conn.RemoteAddr()}),
		checkVer: checkVer,
	}
	c.ctx, c.stop = context.WithCancel(context.Background())

	tc := c.Conn.(*net.TCPConn)
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(1 * time.Minute) // 1min
	return &Conn{smsConn: c, UUID: uuid.New().String()}
}
func (c *smgpConn) Ver() codec.Version {
	return c.Typ
}
func (c *smgpConn) Auth(uid string, pwd string) error {
	// Login to the server.
	req := smgp.NewLoginReq(c.Typ).(*smgp.LoginReq)
	req.ClientID = uid
	req.Secret = pwd
	req.Version = c.Typ

	err := c.SendPDU(req)
	if err != nil {
		c.logger.Errorf("smgp.auth send error: %v", err)
		return err
	}
	p, err := c.RecvPDU()
	if err != nil {
		c.logger.Errorf("smgp.auth recv error: %v", err)
		return err
	}
	var status uint8

	if rsp, ok := p.(*smgp.LoginResp); ok {
		if c.checkVer && rsp.Version != c.Typ {
			return smserror.ErrVersionNotMatch
		}
		status = uint8(rsp.Status)
	} else {
		return smserror.ErrRespNotMatch
	}

	if status != 0 {
		if status <= smserror.ErrnoConnOthers { //ErrnoConnOthers = 5
			err = smserror.ConnRspStatusErrMap[status]
		} else {
			err = smserror.ConnRspStatusErrMap[smserror.ErrnoConnOthers]
		}
		return err
	}
	c.SetState(enum.CONN_AUTHOK)
	return nil
}
func (c *smgpConn) Close() {
	if c != nil {
		if c.State == enum.CONN_CLOSED {
			return
		}
		c.Conn.Close() // close the underlying net.Conn
		c.State = enum.CONN_CLOSED
		c.stop()
	}
}

func (c *smgpConn) SetState(state enum.State) {
	c.State = state
}

// SendPkt pack the smgp packet structure and send it to the other peer.
func (c *smgpConn) SendPDU(pdu codec.PDU) error {
	if c.State == enum.CONN_CLOSED {
		return smserror.ErrConnIsClosed
	}
	if pdu == nil {
		return smserror.ErrPktIsNil
	}
	c.Logger().Debugf("send pdu:%T , %d", pdu, c.Typ)

	buf := codec.NewWriter()
	pdu.Marshal(buf)
	_, err := c.Conn.Write(buf.Bytes()) //block write
	return err
}

func (c *smgpConn) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
func (c *smgpConn) Logger() *logrus.Entry {
	return c.logger
}
func (c *smgpConn) SetReadDeadline(timeout time.Duration) {
	if timeout > 0 {
		c.Conn.SetReadDeadline(time.Now().Add(timeout))
	}
}

// RecvAndUnpackPkt receives smgp byte stream, and unpack it to some smgp packet structure.
func (c *smgpConn) RecvPDU() (codec.PDU, error) {
	if c.State == enum.CONN_CLOSED {
		return nil, smserror.ErrConnIsClosed
	}

	pdu, err := smgp.Parse(c.Conn, c.Typ)
	if err != nil {
		return nil, err
	}

	switch p := pdu.(type) {
	case *smgp.ActiveTestReq: // 当收到心跳请求,内部直接回复心跳,并递归继续获取数据
		resp := p.GetResponse()
		c.SendPDU(resp)
	case *smgp.ActiveTestResp: // 当收到心跳回复,内部直接处理,并递归继续获取数据
		atomic.AddInt32(&c.counter, -1)
	case *smgp.LoginResp: // 当收到登录回复,内部先校验版本
		if c.checkVer && p.Version != c.Typ {
			return nil, fmt.Errorf("smgp version not match [ local: %d != remote: %d ]", c.Typ, p.Version)
		}
	case *smgp.LoginReq: // 当收到登录回复,内部先校验版本
		// 服务端自适应版本
		c.Typ = p.Version
		c.logger = logrus.WithFields(logrus.Fields{"r": c.RemoteAddr(), "v": c.Typ})
	}
	return pdu, nil

}

type smgpListener struct {
	net.Listener
}

func newSmgpListener(l net.Listener) *smgpListener {
	return &smgpListener{l}
}

func (l *smgpListener) accept() (*Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	conn := newSmgpConn(c, smgp.V30, false)
	conn.SetState(enum.CONN_CONNECTED)
	conn.smsConn.(*smgpConn).startActiveTest()
	return conn, nil
}

func (c *smgpConn) startActiveTest() {
	go func() {
		t := time.NewTicker(30 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-c.ctx.Done():
				// once conn close, the goroutine should exit
				return
			case <-t.C:
				// send a active test packet to peer, increase the active test counter
				p := smgp.NewActiveTestReq(c.Typ)
				err := c.SendPDU(p)
				if err != nil {
					c.logger.Errorf("smgp.active send error: %v", err)
				} else {
					atomic.AddInt32(&c.counter, 1)
				}
			}
		}
	}()
}
func (l *smgpListener) Close() error {
	return l.Listener.Close()
}
