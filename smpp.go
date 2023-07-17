package zysms

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/enum"
	"github.com/zhiyin2021/zysms/smpp"
	"github.com/zhiyin2021/zysms/smserror"
)

type smppConn struct {
	net.Conn
	State enum.State
	Typ   codec.Version
	// for SeqId generator goroutine
	// SeqId  <-chan uint32
	// done   chan<- struct{}
	_seqId   uint32
	stop     func()
	ctx      context.Context
	counter  int32
	logger   *logrus.Entry
	checkVer bool
}

// New returns an abstract structure for successfully
// established underlying net.Conn.
func newSmppConn(conn net.Conn, typ codec.Version, checkVer bool) *Conn {
	c := &smppConn{
		Conn:     conn,
		Typ:      typ,
		_seqId:   0,
		logger:   logrus.WithFields(logrus.Fields{"r": conn.RemoteAddr()}),
		checkVer: checkVer,
	}
	c.ctx, c.stop = context.WithCancel(context.Background())

	tc := c.Conn.(*net.TCPConn)
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(1 * time.Minute) // 1min
	return &Conn{smsConn: c, UUID: uuid.New().String()}
}
func (c *smppConn) Ver() codec.Version {
	return c.Typ
}
func (c *smppConn) Auth(uid string, pwd string) error {
	// Login to the server.
	req := smpp.NewBindRequest(smpp.Transceiver)
	req.SystemID = uid
	req.Password = pwd
	req.InterfaceVersion = c.Typ

	err := c.SendPDU(req)
	if err != nil {
		c.logger.Errorf("smpp.auth send error: %v", err)
		return err
	}
	pdu, err := c.RecvPDU()
	if err != nil {
		c.logger.Errorf("smpp.auth recv error: %v", err)
		return err
	}
	if header, ok := pdu.GetHeader().(*smpp.Header); ok {
		status := header.CommandStatus
		if status != smpp.ESME_ROK {
			return fmt.Errorf("login error: %v", status)
		}
		c.SetState(enum.CONN_AUTHOK)
	}
	return nil
}
func (c *smppConn) Close() {
	if c != nil {
		if c.State == enum.CONN_CLOSED {
			return
		}
		c.Conn.Close() // close the underlying net.Conn
		c.State = enum.CONN_CLOSED
		c.stop()
	}
}

func (c *smppConn) SetState(state enum.State) {
	c.State = state
}

// SendPkt pack the smpp packet structure and send it to the other peer.
func (c *smppConn) SendPDU(pdu codec.PDU) error {
	if c.State == enum.CONN_CLOSED {
		return smserror.ErrConnIsClosed
	}
	if pdu == nil {
		return smserror.ErrPktIsNil
	}
	c.Logger().Infof("send pdu:%T , %s", pdu, c.Typ)
	buf := codec.NewWriter()
	pdu.Marshal(buf)
	_, err := c.Conn.Write(buf.Bytes()) //block write

	return err
}

func (c *smppConn) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
func (c *smppConn) Logger() *logrus.Entry {
	return c.logger
}

// RecvAndUnpackPkt receives smpp byte stream, and unpack it to some smpp packet structure.
func (c *smppConn) RecvPDU() (codec.PDU, error) {
	if c.State == enum.CONN_CLOSED {
		return nil, smserror.ErrConnIsClosed
	}
	pdu, err := smpp.Parse(c.Conn)
	if err != nil {
		return nil, err
	}
	if header, ok := pdu.GetHeader().(*smpp.Header); ok {
		switch header.CommandID {
		case smpp.BIND_RECEIVER_RESP, smpp.BIND_TRANSMITTER_RESP, smpp.BIND_TRANSCEIVER_RESP: // 当收到登录回复,内部先校验版本
			if header.CommandStatus != smpp.ESME_ROK {
				return nil, fmt.Errorf("login error: %v", header.CommandStatus)
			}
		case smpp.BIND_RECEIVER, smpp.BIND_TRANSMITTER, smpp.BIND_TRANSCEIVER: /// 当收到登录回复,内部先校验版本
			if v, ok := pdu.(*smpp.BindRequest); ok {
				// 服务端自适应版本
				c.Typ = v.InterfaceVersion
				c.logger = logrus.WithFields(logrus.Fields{"r": c.RemoteAddr(), "v": c.Typ})
			} else {
				return nil, smserror.ErrVersionNotMatch
			}
		}
		return pdu, nil
	}
	return nil, smserror.ErrRespNotMatch
}

type smppListener struct {
	net.Listener
}

func newSmppListener(l net.Listener) *smppListener {
	return &smppListener{l}
}

func (l *smppListener) accept() (*Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	conn := newCmppConn(c, codec.Version(smpp.V34), false)
	conn.SetState(enum.CONN_CONNECTED)
	return conn, nil
}

func (l *smppListener) Close() error {
	return l.Listener.Close()
}
