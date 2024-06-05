package zysms

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"net"
	"runtime/debug"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms/codec"
)

// errors for cmpp server
type (
	PDU codec.PDU
	SMS struct {
		proto        codec.SmsProto
		OnConnect    func(Conn)
		OnDisconnect func(Conn)
		OnError      func(Conn, error)
		OnRecv       func(Conn, PDU)
		// 心跳未响应次数
		OnHeartbeatNoResp func(Conn, int)
		extParam          map[string]string
	}

	Conn interface {
		Close()
		Auth(uid string, pwd string) error
		RemoteAddr() net.Addr
		LocalAddr() net.Addr
		// Recv() ([]byte, error)
		// RecvPDU() (codec.PDU, error)
		SendPDU(PDU) error
		Logger() *logrus.Entry
		Ver() codec.Version
		sendActiveTest() (int32, error)

		SetExtParam(map[string]string)
		GetData() any
		SetData(any)
		SID() string
		Delay() []int64
		IsConnected() bool
		EnabledActiveTest()
	}
)

func New(proto codec.SmsProto) *SMS {
	return &SMS{proto: proto, extParam: map[string]string{}}
}

func (s *SMS) Listen(addr string) (*Listener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	l, err := newListener(ln, s)
	if err != nil {
		return nil, err
	}
	tryGO(func() {
		for {
			sConn, err := l.accept()
			if err != nil {
				logrus.Errorf("listen.accept error:%s", err)
				if e, ok := err.(*net.OpError); ok && e.Error() == "use of closed network connection" {
					return
				}
				continue
			}
			s.run(sConn)

		}
	})
	return l, nil
}
func (s *SMS) ListenTls(addr string, cert []byte, key []byte) (*Listener, error) {
	crt, err := tls.X509KeyPair(cert, key)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	tlsConfig.Certificates = []tls.Certificate{crt}
	// Time returns the current time as the number of seconds since the epoch.
	// If Time is nil, TLS uses time.Now.
	tlsConfig.Time = time.Now
	tlsConfig.Rand = rand.Reader
	ln, err := tls.Listen("tcp", addr, tlsConfig)
	if err != nil {
		return nil, err
	}
	l, err := newListener(ln, s)
	if err != nil {
		return nil, err
	}
	tryGO(func() {
		for {
			sConn, err := l.accept()
			if err != nil {
				logrus.Errorf("listen.accept error:%s", err)
				if e, ok := err.(*net.OpError); ok && e.Error() == "use of closed network connection" {
					return
				}
				continue
			}
			s.run(sConn)
		}
	})
	return l, nil
}
func (s *SMS) doError(conn Conn, err error) {
	if s.OnError != nil {
		if !strings.Contains(err.Error(), "use of closed network connection") {
			s.OnError(conn, err)
		}
	}
}
func (s *SMS) Dial(addr string, uid, pwd string, timeout time.Duration, ext map[string]string) (Conn, error) {
	var err error
	var conn net.Conn
	if ext["tls"] == "1" {
		conn, err = tls.Dial("tcp", addr, &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return nil, err
		}
	} else {
		conn, err = net.DialTimeout("tcp", addr, timeout)
		if err != nil {
			return nil, err
		}
		tc := conn.(*net.TCPConn)
		tc.SetKeepAlive(true)
		tc.SetKeepAlivePeriod(30 * time.Second) // 1min
	}

	sConn := newConn(conn, s)
	if sConn == nil {
		return nil, fmt.Errorf("不支持的协议版本")
	}
	sConn.SetExtParam(ext)
	err = sConn.Auth(uid, pwd)
	if err != nil {
		return nil, err
	}
	// sConn.startActiveTest(s.doError, s.OnHeartbeatNoResp)
	s.run(sConn)
	return sConn, nil
}

func (s *SMS) run(conn *sms_conn) {
	tryGO(func() {
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
			pkt, err := conn.action.recv()
			if err != nil {
				s.doError(conn, err)
				return
			}

			if s.OnRecv != nil {
				// p := &Packet{conn, pkt, nil}
				s.OnRecv(conn, pkt)
				// if !isLogin {
				// 	switch pkt.(type) {
				// 	case *cmpp.ConnReq, *smpp.BindRequest, *smgp.LoginReq, *sgip.BindReq:
				// 		isLogin = true
				// 		conn.startActiveTest(s.doError, s.OnHeartbeatNoResp)
				// 	}
				// }
			}
		}
	})
}

type Listener struct {
	net.Listener
	parent *SMS
	// extParam map[string]string
	// proto    codec.SmsProto
}

func newListener(l net.Listener, parent *SMS) (*Listener, error) {
	switch parent.proto {
	case codec.CMPP20, codec.CMPP21, codec.CMPP30, codec.SMGP30, codec.SGIP, codec.SMPP33, codec.SMPP34:
	default:
		return nil, fmt.Errorf("不支持的协议版本")
	}
	return &Listener{l, parent}, nil
}

func (l *Listener) accept() (*sms_conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	tc := c.(*net.TCPConn)
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(30 * time.Second) // 1min

	conn := newConn(c, l.parent)
	if conn == nil {
		return nil, fmt.Errorf("不支持的协议版本")
	}
	return conn, nil
}

func (l *Listener) Close() error {
	return l.Listener.Close()
}

func tryGO(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logrus.StandardLogger().Fatalf("panic:%v\n%s", err, debug.Stack())
			}
		}()
		f()
	}()
}
