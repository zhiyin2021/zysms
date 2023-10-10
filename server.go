package zysms

import (
	"context"
	"net"
	"strconv"
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
	SMS struct {
		proto        codec.SmsProto
		OnConnect    func(*Conn)
		OnDisconnect func(*Conn)
		OnError      func(*Conn, error)
		OnRecv       func(*Packet) error
		// NodeId       uint32 // sgip 序列号使用
		activeCount    int32
		activeInterval int32
		extParam       map[string]string
	}

	// smsOption struct {
	// 	activeInterval  time.Duration
	// 	activeFailCount int32
	// 	extParam        map[string]string
	// }
	// Opt func(*smsOption)

	Conn struct {
		smsConn
		Data any
		// Logger *logrus.Entry
		UUID string
		ctx  context.Context
		stop func()
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
		sendActiveTest() (int32, error)
	}
)

func New(proto codec.SmsProto, extParam map[string]string) *SMS {
	// smsOpt := smsOption{activeInterval: 5 * time.Second, activeFailCount: 3, extParam: map[string]string{}}
	// for _, opt := range opts {
	// 	opt(&smsOpt)
	// }
	activeCount := 3
	activeInterval := 5
	if extParam != nil {
		if extParam["active_count"] != "" {
			n, err := strconv.Atoi(extParam["active_count"])
			if err == nil {
				activeCount = n
			}
		}
		if extParam["active_interval"] != "" {
			n, err := strconv.Atoi(extParam["active_interval"])
			if err == nil {
				activeInterval = n
			}
		}
	}
	return &SMS{proto: proto, extParam: extParam, activeCount: int32(activeCount), activeInterval: int32(activeInterval)}
}

func (s *SMS) Listen(addr string) (smsListener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	var l smsListener
	switch s.proto {
	case codec.CMPP20, codec.CMPP21, codec.CMPP30:
		l = newCmppListener(ln, s.extParam)
	case codec.SMPP33, codec.SMPP34:
		l = newSmppListener(ln, s.extParam)
	case codec.SMGP13, codec.SMGP20, codec.SMGP30:
		l = newSmgpListener(ln, s.extParam)
	case codec.SGIP:
		l = newSgipListener(ln, s.extParam)
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

func (s *SMS) Dial(addr string, uid, pwd string, timeout time.Duration, checkVer bool) (*Conn, error) {
	var err error
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, err
	}
	var zconn *Conn
	switch s.proto {
	case codec.CMPP20:
		zconn = newCmppConn(conn, cmpp.V20, checkVer, s.extParam)
	case codec.CMPP21:
		zconn = newCmppConn(conn, cmpp.V21, checkVer, s.extParam)
	case codec.CMPP30:
		zconn = newCmppConn(conn, cmpp.V30, checkVer, s.extParam)
	case codec.SMPP33:
		zconn = newSmppConn(conn, smpp.V33, checkVer, s.extParam)
	case codec.SMPP34:
		zconn = newSmppConn(conn, smpp.V34, checkVer, s.extParam)
	case codec.SMGP30:
		zconn = newSmgpConn(conn, smgp.V30, checkVer, s.extParam)
	case codec.SGIP:
		zconn = newSgipConn(conn, sgip.V12, checkVer, s.extParam)
	default:
		return nil, smserror.ErrProtoNotSupport
	}
	zconn.ctx, zconn.stop = context.WithCancel(context.Background())
	zconn.SetState(enum.CONN_CONNECTED)
	err = zconn.Auth(uid, pwd)
	if err != nil {
		return nil, err
	}
	zconn.startActiveTest(s.activeInterval, s.activeCount)
	go s.run(zconn)
	return zconn, nil
}

func (s *SMS) run(conn *Conn) {
	if s.OnConnect != nil {
		s.OnConnect(conn)
	}
	defer func() {
		if s.OnDisconnect != nil {
			s.OnDisconnect(conn)
		}
		conn.Close()
	}()
	conn.startActiveTest(s.activeInterval, s.activeCount)
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

func (c *Conn) startActiveTest(activeInterval, activeCount int32) {

	go func() {
		t := time.NewTicker(time.Duration(activeInterval) * time.Second)
		defer t.Stop()
		for {
			select {
			case <-c.ctx.Done():
				// once conn close, the goroutine should exit
				return
			case <-t.C:
				n, err := c.sendActiveTest()
				if err != nil {
					c.Logger().Errorln(err)
					return
				}
				if n > activeCount {
					c.Logger().Errorf("超过3次心跳未收到响应,关闭连接")
					c.Close()
					return
				}
			}
		}
	}()
}

// func WithActiveInterval(interval time.Duration) Opt {
// 	return func(opt *smsOption) {
// 		opt.activeInterval = interval
// 	}
// }
// func WithActiveFailCount(count int32) Opt {
// 	return func(opt *smsOption) {
// 		opt.activeFailCount = count
// 	}
// }
// func WithExtParam(extParam map[string]string) Opt {
// 	return func(opt *smsOption) {
// 		opt.extParam = extParam
// 	}
// }
