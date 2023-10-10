package zysms

import (
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
	_seqId     uint32
	counter    int32
	logger     *logrus.Entry
	checkVer   bool
	OnError    func(*Conn, error)
	activeFail int32
	extParam   map[string]string
	// activePeer bool // 默认false,当前连接发送心跳请求, 当收到对方心跳请求后,设置true,不再发送心跳请求
	// activeLast time.Time
}

// New returns an abstract structure for successfully
// established underlying net.Conn.
func newSmppConn(conn net.Conn, typ codec.Version, checkVer bool, extParam map[string]string) *Conn {
	c := &smppConn{
		Conn:     conn,
		Typ:      typ,
		_seqId:   0,
		logger:   logrus.WithFields(logrus.Fields{"r": conn.RemoteAddr()}),
		checkVer: checkVer,
		extParam: extParam,
	}
	tc := c.Conn.(*net.TCPConn)
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(10 * time.Second) // 1min
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
	if c.extParam != nil && c.extParam["system_type"] != "" {
		req.SystemType = c.extParam["system_type"]
	}
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
	}
}

func (c *smppConn) SetState(state enum.State) {
	c.State = state
}

// SendPkt pack the smpp packet structure and send it to the other peer.
func (c *smppConn) SendPDU(pdu codec.PDU) error {
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
		// if !c.activePeer {
		// 	c.activePeer = true
		// }
	case *smpp.EnquireLinkResp: // 当收到心跳回复,内部直接处理,并递归继续获取数据
		atomic.AddInt32(&c.counter, -1)
		// c.activeLast = time.Now()
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

func (c *smppConn) sendActiveTest() (int32, error) {
	p := smpp.NewEnquireLink()
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

type smppListener struct {
	net.Listener
	extParam map[string]string
}

func newSmppListener(l net.Listener, extParam map[string]string) *smppListener {
	return &smppListener{l, extParam}
}

func (l *smppListener) accept() (*Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	conn := newSmppConn(c, codec.Version(smpp.V34), false, l.extParam)
	conn.SetState(enum.CONN_CONNECTED)
	return conn, nil
}

func (l *smppListener) Close() error {
	return l.Listener.Close()
}
