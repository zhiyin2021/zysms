package zysms

import (
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms/enum"
	"github.com/zhiyin2021/zysms/proto"
)

type SmsConn interface {
	Close()
	Auth(uid string, pwd string, timeout time.Duration) error
	RemoteAddr() net.Addr
	// Recv() ([]byte, error)
	RecvPkt(time.Duration) (proto.Packer, error)
	SendPkt(proto.Packer, uint32) error
	SetState(enum.State)
	Logger() *logrus.Entry
}
