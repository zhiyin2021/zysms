package zysms

import (
	"context"
	"net"
	"sync/atomic"
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
	c.startActiveTest()

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
	req.SystemType = "NULL"
	err := c.SendPDU(req)
	if err != nil {
		return err
	}
	pdu, err := c.RecvPDU()
	if err != nil {
		return err
	}
	if header, ok := pdu.GetHeader().(*smpp.Header); ok {
		status := header.CommandStatus
		if status != smpp.ESME_ROK {
			return smserror.NewSmsErr(int(status), "smpp.login.error") //fmt.Errorf("login error: %v", status)
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
	c.Logger().Debugf("send pdu:%T , %d", pdu, c.Typ)
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
	switch p := pdu.(type) {
	case *smpp.EnquireLink: // 当收到心跳请求,内部直接回复心跳,并递归继续获取数据
		resp := p.GetResponse()
		c.SendPDU(resp)
	case *smpp.EnquireLinkResp: // 当收到心跳回复,内部直接处理,并递归继续获取数据
		atomic.AddInt32(&c.counter, -1)
	case *smpp.BindResp: // 当收到登录回复,内部先校验版本
		if p.CommandStatus != smpp.ESME_ROK {
			return nil, smserror.NewSmsErr(int(p.CommandStatus), "smpp.login.error")
		}
	case *smpp.BindRequest: /// 当收到登录回复,内部先校验版本
		// 服务端自适应版本
		c.Typ = p.InterfaceVersion
		c.logger = logrus.WithFields(logrus.Fields{"r": c.RemoteAddr(), "v": c.Typ})
	}
	return pdu, nil
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
	conn := newSmppConn(c, codec.Version(smpp.V34), false)
	conn.SetState(enum.CONN_CONNECTED)

	return conn, nil
}

func (c *smppConn) startActiveTest() {
	go func() {
		fail := 0
		t := time.NewTicker(30 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-c.ctx.Done():
				// once conn close, the goroutine should exit
				return
			case <-t.C:
				// send a active test packet to peer, increase the active test counter
				p := smpp.NewEnquireLink()
				err := c.SendPDU(p)
				if err != nil {
					fail++
					c.logger.Errorf("smpp.active send error: %v", err)
					if fail > 3 {
						return
					}
				} else {
					fail = 0
					atomic.AddInt32(&c.counter, 1)
				}
			}
		}
	}()
}

func (l *smppListener) Close() error {
	return l.Listener.Close()
}
