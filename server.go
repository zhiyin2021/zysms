package zysms

import (
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms/cmpp"
	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/enum"
	"github.com/zhiyin2021/zysms/sgip"
	"github.com/zhiyin2021/zysms/smgp"
	"github.com/zhiyin2021/zysms/smpp"
	"github.com/zhiyin2021/zysms/smserror"
)

// errors for cmpp server
type (
	Packet struct {
		Conn *Conn
		Req  codec.PDU
		Resp codec.PDU
	}
	// handleEvent func(*Conn, codec.Packer) error
	sms struct {
		proto        codec.SmsProto
		OnConnect    func(*Conn)
		OnDisconnect func(*Conn)
		OnError      func(*Conn, error)
		OnRecv       func(*Packet) error
		NodeId       uint32 // sgip 序列号使用
	}
	Conn struct {
		smsConn
		// Data   any
		// Logger *logrus.Entry
		UUID string
	}
	smsListener interface {
		accept() (*Conn, error)
		Close() error
	}

	smsConn interface {
		Close()
		Auth(uid string, pwd string) error
		RemoteAddr() net.Addr
		// Recv() ([]byte, error)
		RecvPDU() (codec.PDU, error)
		SendPDU(codec.PDU) error
		SetState(enum.State)
		Logger() *logrus.Entry
		Ver() codec.Version
	}
)

func New(proto codec.SmsProto) *sms {
	return &sms{proto: proto}
}

func (s *sms) Listen(addr string) (smsListener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	var l smsListener
	switch s.proto {
	case codec.CMPP20, codec.CMPP21, codec.CMPP30:
		l = newCmppListener(ln)
	case codec.SMPP33, codec.SMPP34:
		l = newSmppListener(ln)
	case codec.SMGP13, codec.SMGP20, codec.SMGP30:
		l = newSmgpListener(ln)
	case codec.SGIP:
		l = newSgipListener(ln, s.NodeId)
	}
	go func() {
		for {
			conn, err := l.accept()
			if err != nil {
				logrus.Errorf("listen.accept error:%s", err)
				if e, ok := err.(*net.OpError); ok && e.Error() == "use of closed network connection" {
					return
				}
				continue
			}
			go s.run(conn)
		}
	}()
	return l, nil
}

func (s *sms) Dial(addr string, uid, pwd string, timeout time.Duration, checkVer bool) (*Conn, error) {
	var err error
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, err
	}
	var zconn *Conn
	switch s.proto {
	case codec.CMPP20:
		zconn = newCmppConn(conn, cmpp.V20, checkVer)
	case codec.CMPP21:
		zconn = newCmppConn(conn, cmpp.V21, checkVer)
	case codec.CMPP30:
		zconn = newCmppConn(conn, cmpp.V30, checkVer)
	case codec.SMPP33:
		zconn = newSmppConn(conn, smpp.V33, checkVer)
	case codec.SMPP34:
		zconn = newSmppConn(conn, smpp.V34, checkVer)
	case codec.SMGP30:
		zconn = newSmgpConn(conn, smgp.V30, checkVer)
	case codec.SGIP:
		zconn = newSgipConn(conn, sgip.V12, checkVer, s.NodeId)
	default:
		return nil, smserror.ErrProtoNotSupport
	}
	zconn.SetState(enum.CONN_CONNECTED)
	err = zconn.Auth(uid, pwd)
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
		pkt, err := conn.RecvPDU()
		if err != nil {
			if s.OnError != nil {
				s.OnError(conn, err)
			}
			return
		}
		if s.OnRecv != nil {
			p := &Packet{conn, pkt, nil}
			err = s.OnRecv(p)
			if p.Resp != nil {
				err := conn.SendPDU(p.Resp)
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
