package zysms

import (
	"fmt"
	"net"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms/cmpp"
	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/enum"
	"github.com/zhiyin2021/zysms/sgip"
	"github.com/zhiyin2021/zysms/smgp"
	"github.com/zhiyin2021/zysms/smpp"
)

// errors for cmpp server
type (
	PDU codec.PDU
	SMS struct {
		proto        codec.SmsProto
		OnConnect    func(Conn)
		OnDisconnect func(Conn)
		OnError      func(Conn, error)
		OnRecv       func(Conn, PDU) (PDU, error)
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
		SetState(enum.State)
		Logger() *logrus.Entry
		Ver() codec.Version
		sendActiveTest() (int32, error)

		SetExtParam(map[string]string)
		GetData() any
		SetData(any)
		UID() string
	}
)

func New(proto codec.SmsProto) *SMS {
	// smsOpt := smsOption{activeInterval: 5 * time.Second, activeFailCount: 3, extParam: map[string]string{}}
	// for _, opt := range opts {
	// 	opt(&smsOpt)
	// }
	return &SMS{proto: proto, extParam: map[string]string{}}
}

func (s *SMS) Listen(addr string) (*Listener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	l, err := newListener(ln, s.proto, s.extParam)
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
			s.run(sConn, false)

		}
	})
	return l, nil
}

func (s *SMS) Dial(addr string, uid, pwd string, timeout time.Duration, ext map[string]string) (Conn, error) {
	var err error
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, err
	}
	tc := conn.(*net.TCPConn)
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(30 * time.Second) // 1min

	sConn := newConn(conn, s.proto)
	if sConn == nil {
		return nil, fmt.Errorf("不支持的协议版本")
	}
	sConn.SetExtParam(ext)
	err = sConn.Auth(uid, pwd)
	if err != nil {
		return nil, err
	}
	sConn.startActiveTest(s.OnError, s.OnHeartbeatNoResp)
	s.run(sConn, true)
	return sConn, nil
}

func (s *SMS) run(conn *sms_conn, isLogin bool) {
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
				if s.OnError != nil {
					s.OnError(conn, err)
				}
				return
			}

			if s.OnRecv != nil {
				// p := &Packet{conn, pkt, nil}
				resp, err := s.OnRecv(conn, pkt)
				if !isLogin {
					switch pkt.(type) {
					case *cmpp.ConnReq, *smpp.BindRequest, *smgp.LoginReq, *sgip.BindReq:
						isLogin = true
						conn.startActiveTest(s.OnError, s.OnHeartbeatNoResp)
					}
				}
				if resp != nil {
					err := conn.SendPDU(resp)
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
	})
}

type Listener struct {
	net.Listener
	extParam map[string]string
	proto    codec.SmsProto
}

func newListener(l net.Listener, proto codec.SmsProto, extParam map[string]string) (*Listener, error) {
	switch proto {
	case codec.CMPP20, codec.CMPP21, codec.CMPP30, codec.SMGP30, codec.SGIP, codec.SMPP33, codec.SMPP34:
	default:
		return nil, fmt.Errorf("不支持的协议版本")
	}
	return &Listener{l, extParam, proto}, nil
}

func (l *Listener) accept() (*sms_conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	tc := c.(*net.TCPConn)
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(30 * time.Second) // 1min

	conn := newConn(c, l.proto)
	if conn == nil {
		return nil, fmt.Errorf("不支持的协议版本")
	}
	conn.SetState(enum.CONN_CONNECTED)
	return conn, nil
}

func (l *Listener) Close() error {
	return l.Listener.Close()
}

func tryGO(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 1<<16)
				runtime.Stack(buf, true)
				logrus.Errorf("panic:%v\n%s", err, buf)
			}
		}()
		f()
	}()
}
