package zysms

import (
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms/cmpp"
	"github.com/zhiyin2021/zysms/enum"
	"github.com/zhiyin2021/zysms/event"
	"github.com/zhiyin2021/zysms/proto"
	"github.com/zhiyin2021/zysms/smserror"
)

// errors for cmpp server
type (
	handleEvent func(*Conn, proto.Packer) error
	sms         struct {
		server       smsListener
		client       smsConn
		proto        proto.SmsProto
		events       map[event.SmsEvent]handleEvent
		OnConnect    func(*Conn)
		OnDisconnect func(*Conn)
		OnError      func(*Conn, error)
	}
	Conn struct {
		smsConn
		Data   any
		Logger *logrus.Entry
	}
	smsListener interface {
		Accept() (*Conn, error)
		Close() error
	}

	smsConn interface {
		Close()
		Auth(uid string, pwd string, timeout time.Duration) error
		RemoteAddr() net.Addr
		// Recv() ([]byte, error)
		RecvPkt(time.Duration) (proto.Packer, error)
		SendPkt(proto.Packer, uint32) error
		SetState(enum.State)
	}
)

func New(proto proto.SmsProto) *sms {
	return &sms{proto: proto, events: make(map[event.SmsEvent]handleEvent)}
}

func (s *sms) Listen(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	var l smsListener
	switch s.proto {
	case proto.CMPP2:
		l = newCmppListener(ln, cmpp.V20)
	case proto.CMPP3:
		l = newCmppListener(ln, cmpp.V30)
	}
	s.server = l
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go s.run(conn)
	}
}

func (s *sms) Dial(addr string, uid, pwd string, timeout time.Duration) (*Conn, error) {
	var err error
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {

		return nil, err
	}
	var zconn *Conn
	switch s.proto {
	case proto.CMPP2:
		zconn = newCmppConn(conn, cmpp.V20)
	case proto.CMPP3:
		zconn = newCmppConn(conn, cmpp.V30)
	default:
		return nil, smserror.ErrProtoNotSupport
	}
	zconn.SetState(enum.CONN_CONNECTED)
	err = zconn.Auth(uid, pwd, timeout)
	if err != nil {
		return nil, err
	}
	go s.run(zconn)
	return zconn, nil
}

func (s *sms) Handle(e event.SmsEvent, f handleEvent) {
	s.events[e] = f
}

func (s *sms) run(conn *Conn) {
	if s.OnConnect != nil {
		s.OnConnect(conn)
	}
	defer func() {
		if s.OnDisconnect != nil {
			s.OnDisconnect(conn)
		}
		conn.Close()
	}()

	for {
		pkt, err := conn.RecvPkt(0)
		if err != nil {
			s.doError(conn, err)
			return
		}
		if e, ok := s.events[pkt.Event()]; ok {
			err = e(conn, pkt)
			if err != nil {
				return
			}
		} else if e, ok := s.events[event.SmsEventUnknown]; ok {
			err = e(conn, pkt)
			if err != nil {
				return
			}
		}
	}
}
func (s *sms) doError(c *Conn, e error) {
	if s.OnError != nil {
		s.OnError(c, e)
	}
}
