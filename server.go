package zysms

import (
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms/cmpp"
	"github.com/zhiyin2021/zysms/enum"
	"github.com/zhiyin2021/zysms/proto"
	"github.com/zhiyin2021/zysms/smserror"
)

// errors for cmpp server
type (
	Packet struct {
		Conn *Conn
		Req  proto.Packer
		Resp proto.Packer
	}
	// handleEvent func(*Conn, proto.Packer) error
	sms struct {
		proto        proto.SmsProto
		OnConnect    func(*Conn)
		OnDisconnect func(*Conn)
		OnError      func(*Conn, error)
		OnEvent      func(*Packet) error
	}
	Conn struct {
		smsConn
		// Data   any
		Logger *logrus.Entry
		Proto  proto.SmsProto
		UUID   string
	}
	smsListener interface {
		accept() (*Conn, error)
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
	return &sms{proto: proto}
}

func (s *sms) Listen(addr string) (smsListener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	var l smsListener
	switch s.proto {
	case proto.CMPP2:
		l = newCmppListener(ln, cmpp.V20)
	case proto.CMPP3:
		l = newCmppListener(ln, cmpp.V30)
	}
	go func() {
		for {
			conn, err := l.accept()
			if err != nil {
				return
			}
			conn.Proto = s.proto
			go s.run(conn)
		}
	}()
	return l, nil
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
	zconn.Proto = s.proto
	zconn.SetState(enum.CONN_CONNECTED)
	err = zconn.Auth(uid, pwd, timeout)
	if err != nil {
		return nil, err
	}
	go s.run(zconn)
	return zconn, nil
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
			if s.OnError != nil {
				s.OnError(conn, err)
			}
			return
		}
		if s.OnEvent != nil {
			p := &Packet{conn, pkt, nil}
			err = s.OnEvent(p)
			if p.Resp != nil {
				err := conn.SendPkt(p.Resp, pkt.SeqId())
				if err != nil {
					if s.OnError != nil {
						s.OnError(conn, err)
					}
					return
				}
			}
			if err != nil {
				return
			}
		}
	}
}
