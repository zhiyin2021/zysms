package zysms

import (
	"net"
	"time"

	"github.com/zhiyin2021/zysms/cmpp"
	"github.com/zhiyin2021/zysms/enum"
)

// errors for cmpp server
type SmsProto byte

const (
	CMPP2 SmsProto = iota
	CMPP3
	SMGP
	SGIP
	SMPP
)

type SmsListener interface {
	Accept() (SmsConn, error)
	Close() error
}

func Listen(addr string, proto SmsProto) (SmsListener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	var l SmsListener
	switch proto {
	case CMPP2:
		l = newCmppListener(ln, cmpp.V20)
	case CMPP3:
		l = newCmppListener(ln, cmpp.V30)
	}
	return l, nil
}
func Dial(addr string, proto SmsProto, uid, pwd string, timeout time.Duration) (SmsConn, error) {
	var err error
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, err
	}
	var smsConn SmsConn
	switch proto {
	case CMPP2:
		smsConn = newCmppConn(conn, cmpp.V20)
	case CMPP3:
		smsConn = newCmppConn(conn, cmpp.V30)
	}
	smsConn.SetState(enum.CONN_CONNECTED)
	err = smsConn.Auth(uid, pwd, timeout)
	return smsConn, err
}
