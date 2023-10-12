package zysms

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms/cmpp"
	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/enum"
	"github.com/zhiyin2021/zysms/sgip"
	"github.com/zhiyin2021/zysms/smgp"
	"github.com/zhiyin2021/zysms/smpp"
	"github.com/zhiyin2021/zysms/smserror"
	"github.com/zhiyin2021/zysms/utils"
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
		extParam     map[string]string
	}

	Conn struct {
		smsConn
		Data any
		// Logger *logrus.Entry
		UUID           string
		ctx            context.Context
		stop           func()
		activeCount    int32
		activeInterval int
	}
	smsListener interface {
		accept() (smsConn, error)
		Close() error
	}

	smsConn interface {
		close()
		Auth(uid string, pwd string) error
		RemoteAddr() net.Addr
		// Recv() ([]byte, error)
		RecvPDU() (codec.PDU, error)
		SendPDU(codec.PDU) error
		SetState(enum.State)
		Logger() *logrus.Entry
		Ver() codec.Version
		sendActiveTest() (int32, error)
		setExtParam(map[string]string)
	}
)

func New(proto codec.SmsProto) *SMS {
	// smsOpt := smsOption{activeInterval: 5 * time.Second, activeFailCount: 3, extParam: map[string]string{}}
	// for _, opt := range opts {
	// 	opt(&smsOpt)
	// }
	return &SMS{proto: proto, extParam: map[string]string{}}
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
			sConn, err := l.accept()
			if err != nil {
				logrus.Errorf("listen.accept error:%s", err)
				if e, ok := err.(*net.OpError); ok && e.Error() == "use of closed network connection" {
					return
				}
				continue
			}
			zconn := &Conn{smsConn: sConn, UUID: uuid.New().String(), activeCount: 3, activeInterval: 5}
			zconn.ctx, zconn.stop = context.WithCancel(context.Background())
			zconn.SetState(enum.CONN_CONNECTED)
			go s.run(zconn, false)
		}
	}()
	return l, nil
}

func (s *SMS) Dial(addr string, uid, pwd string, timeout time.Duration, ext map[string]string) (*Conn, error) {

	var err error
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, err
	}
	tc := conn.(*net.TCPConn)
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(30 * time.Second) // 1min

	var sConn smsConn
	switch s.proto {
	case codec.CMPP20:
		sConn = newCmppConn(conn, cmpp.V20)
	case codec.CMPP21:
		sConn = newCmppConn(conn, cmpp.V21)
	case codec.CMPP30:
		sConn = newCmppConn(conn, cmpp.V30)
	case codec.SMPP33:
		sConn = newSmppConn(conn, smpp.V33)
	case codec.SMPP34:
		sConn = newSmppConn(conn, smpp.V34)
	case codec.SMGP30:
		sConn = newSmgpConn(conn, smgp.V30)
	case codec.SGIP:
		sConn = newSgipConn(conn, sgip.V12)
	default:
		return nil, smserror.ErrProtoNotSupport
	}
	zconn := &Conn{smsConn: sConn, UUID: uuid.New().String(), activeCount: 3, activeInterval: 5}
	zconn.ctx, zconn.stop = context.WithCancel(context.Background())
	zconn.SetState(enum.CONN_CONNECTED)
	zconn.SetExtParam(ext)
	err = zconn.Auth(uid, pwd)
	if err != nil {
		return nil, err
	}
	zconn.startActiveTest(s.OnError)
	go s.run(zconn, true)
	return zconn, nil
}

func (s *SMS) run(conn *Conn, isLogin bool) {
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
			if !isLogin {
				switch pkt.(type) {
				case *cmpp.ConnReq, *smpp.BindRequest, *smgp.LoginReq, *sgip.BindReq:
					isLogin = true
					conn.startActiveTest(s.OnError)
				}
			}
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

func (c *Conn) startActiveTest(errEvent func(*Conn, error)) {
	if c.activeInterval > 0 && c.activeCount > 0 {
		go func() {
			t := time.NewTicker(time.Duration(c.activeInterval) * time.Second)
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
					if n >= c.activeCount {
						if errEvent != nil {
							errEvent(c, fmt.Errorf("%d次心跳间隔(%ds)未收到响应,关闭连接", c.activeCount, c.activeInterval))
						} else {
							c.Logger().Errorf("%d次心跳间隔(%ds)未收到响应,关闭连接", c.activeCount, c.activeInterval)
						}
						c.Close()
						return
					}
				}
			}
		}()
	}
}

func (c *Conn) Close() {
	c.stop()
	c.close()
}

/*
active_count 心跳失败次数
active_interval 心跳间隔
check_version 是否校验版本
system_type 系统类型[smpp 特有]
*/
func (c *Conn) SetExtParam(ext map[string]string) {
	c.setExtParam(ext)
	c.activeCount = utils.MapItem(ext, "active_count", int32(3))
	c.activeInterval = utils.MapItem(ext, "active_interval", int(5))
}
