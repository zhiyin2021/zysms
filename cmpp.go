package zysms

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms/cmpp"
	"github.com/zhiyin2021/zysms/enum"
	"github.com/zhiyin2021/zysms/proto"
	"github.com/zhiyin2021/zysms/smserror"
)

type cmppConn struct {
	net.Conn
	State enum.State
	Typ   cmpp.Version
	// for SeqId generator goroutine
	// SeqId  <-chan uint32
	// done   chan<- struct{}
	_seqId   uint32
	stop     func()
	ctx      context.Context
	counter  int32
	logger   *logrus.Entry
	checkVer bool
}

// New returns an abstract structure for successfully
// established underlying net.Conn.
func newCmppConn(conn net.Conn, typ cmpp.Version, checkVer bool) *Conn {
	c := &cmppConn{
		Conn:     conn,
		Typ:      typ,
		_seqId:   0,
		logger:   logrus.WithFields(logrus.Fields{"r": conn.RemoteAddr()}),
		checkVer: checkVer,
	}
	c.ctx, c.stop = context.WithCancel(context.Background())

	tc := c.Conn.(*net.TCPConn)
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(1 * time.Minute) // 1min
	return &Conn{smsConn: c, UUID: uuid.New().String()}
}
func (c *cmppConn) Proto() proto.SmsProto {
	return c.Typ.Proto()
}
func (c *cmppConn) Auth(uid string, pwd string, timeout time.Duration) error {
	// Login to the server.
	req := &cmpp.CmppConnReq{
		SrcAddr: uid,
		Secret:  pwd,
		Version: c.Typ,
	}
	err := c.SendPkt(req, c.seqId())
	if err != nil {
		c.logger.Errorf("cmpp.auth send error: %v", err)
		return err
	}
	p, err := c.RecvPkt(timeout)
	if err != nil {
		c.logger.Errorf("cmpp.auth recv error: %v", err)
		return err
	}
	var status uint8

	if rsp, ok := p.(*cmpp.CmppConnRsp); ok {
		if c.checkVer && rsp.Version != c.Typ {
			return smserror.ErrVersionNotMatch
		}
		status = uint8(rsp.Status)
	} else {
		return smserror.ErrRespNotMatch
	}

	if status != 0 {
		if status <= cmpp.ErrnoConnOthers { //ErrnoConnOthers = 5
			err = cmpp.ConnRspStatusErrMap[status]
		} else {
			err = cmpp.ConnRspStatusErrMap[cmpp.ErrnoConnOthers]
		}
		return err
	}
	c.SetState(enum.CONN_AUTHOK)
	return nil
}
func (c *cmppConn) Close() {
	if c != nil {
		if c.State == enum.CONN_CLOSED {
			return
		}
		c.Conn.Close() // close the underlying net.Conn
		c.State = enum.CONN_CLOSED
		c.stop()
	}
}

func (c *cmppConn) SetState(state enum.State) {
	c.State = state
}
func (c *cmppConn) seqId() uint32 {
	sid := atomic.AddUint32(&c._seqId, 1)
	return sid
}

// SendPkt pack the cmpp packet structure and send it to the other peer.
func (c *cmppConn) SendPkt(pkt proto.Packer, seqId uint32) error {
	if c.State == enum.CONN_CLOSED {
		return smserror.ErrConnIsClosed
	}
	if pkt == nil {
		return smserror.ErrPktIsNil
	}
	if seqId == 0 {
		seqId = c.seqId()
	}
	c.Logger().Infof("send pkt:%T , %s", pkt, c.Typ)

	data := pkt.Pack(seqId, c.Typ.Proto())

	_, err := c.Conn.Write(data) //block write
	if err != nil {
		return err
	}
	return nil
}

const (
	defaultReadBufferSize = 4096
)

// readBuffer is used to optimize the performance of
// RecvAndUnpackPkt.
type readBuffer struct {
	totalLen  uint32
	commandId cmpp.CommandId
	leftData  [defaultReadBufferSize]byte
}

var readBufferPool = sync.Pool{
	New: func() any {
		return &readBuffer{}
	},
}

func (c *cmppConn) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
func (c *cmppConn) Logger() *logrus.Entry {
	return c.logger
}

// RecvAndUnpackPkt receives cmpp byte stream, and unpack it to some cmpp packet structure.
func (c *cmppConn) RecvPkt(timeout time.Duration) (proto.Packer, error) {
	if c.State == enum.CONN_CLOSED {
		return nil, smserror.ErrConnIsClosed
	}
	rb := readBufferPool.Get().(*readBuffer)
	defer readBufferPool.Put(rb)
	defer c.SetReadDeadline(time.Time{})

	// Total_Length in packet
	if timeout != 0 {
		c.SetReadDeadline(time.Now().Add(timeout))
	}
	err := binary.Read(c.Conn, binary.BigEndian, &rb.totalLen)
	if err != nil {
		return nil, err
	}

	if c.Typ == cmpp.V30 && (rb.totalLen < cmpp.CMPP3_PACKET_MIN || rb.totalLen > cmpp.CMPP3_PACKET_MAX) {
		return nil, smserror.ErrTotalLengthInvalid
	} else if rb.totalLen < cmpp.CMPP2_PACKET_MIN || rb.totalLen > cmpp.CMPP2_PACKET_MAX {
		return nil, smserror.ErrTotalLengthInvalid
	}

	// Command_Id
	if timeout != 0 {
		c.SetReadDeadline(time.Now().Add(timeout))
	}
	err = binary.Read(c.Conn, binary.BigEndian, &rb.commandId)
	if err != nil {
		netErr, ok := err.(net.Error)
		if ok {
			if netErr.Timeout() {
				return nil, smserror.ErrReadCmdIDTimeout
			}
		}
		return nil, err
	}

	if !((rb.commandId > cmpp.CMPP_REQUEST_MIN && rb.commandId < cmpp.CMPP_REQUEST_MAX) ||
		(rb.commandId > cmpp.CMPP_RESPONSE_MIN && rb.commandId < cmpp.CMPP_RESPONSE_MAX)) {
		return nil, smserror.ErrCommandIdInvalid
	}

	// The left packet data (start from seqId in header).
	if timeout != 0 {
		c.SetReadDeadline(time.Now().Add(timeout))
	}
	var leftData = rb.leftData[0:(rb.totalLen - 8)]
	_, err = io.ReadFull(c.Conn, leftData)
	if err != nil {
		return nil, err
	}

	if fun, ok := cmpp.CmppPacket[rb.commandId]; ok {
		p := fun(c.Typ, leftData)
		switch rb.commandId {
		case cmpp.CMPP_ACTIVE_TEST: // 当收到心跳请求,内部直接回复心跳,并递归继续获取数据
			resp := &cmpp.CmppActiveTestRsp{}
			c.SendPkt(resp, p.SeqId())
			return c.RecvPkt(timeout)
		case cmpp.CMPP_ACTIVE_TEST_RESP: // 当收到心跳回复,内部直接处理,并递归继续获取数据
			atomic.AddInt32(&c.counter, -1)
			return c.RecvPkt(timeout)
		case cmpp.CMPP_CONNECT_RESP: // 当收到登录回复,内部先校验版本
			if v, ok := p.(*cmpp.CmppConnRsp); ok {
				if c.checkVer && v.Version != c.Typ {
					return nil, fmt.Errorf("cmpp version not match [ local: %d != remote: %d ]", c.Typ, v.Version)
				}
			}
		case cmpp.CMPP_CONNECT: // 当收到登录回复,内部先校验版本
			if v, ok := p.(*cmpp.CmppConnReq); ok {
				// 服务端自适应版本
				c.Typ = v.Version
				c.logger = logrus.WithFields(logrus.Fields{"r": c.RemoteAddr(), "v": c.Typ})
			} else {
				return nil, smserror.ErrVersionNotMatch
			}
		}
		return p, nil
	}
	return nil, smserror.ErrCommandIdNotSupported

}

type cmppListener struct {
	net.Listener
}

func newCmppListener(l net.Listener) *cmppListener {
	return &cmppListener{l}
}

func (l *cmppListener) accept() (*Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	conn := newCmppConn(c, cmpp.V30, false)
	conn.SetState(enum.CONN_CONNECTED)
	conn.smsConn.(*cmppConn).startActiveTest()
	return conn, nil
}

func (c *cmppConn) startActiveTest() {
	go func() {
		t := time.NewTicker(10 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-c.ctx.Done():
				// once conn close, the goroutine should exit
				return
			case <-t.C:
				// send a active test packet to peer, increase the active test counter
				p := &cmpp.CmppActiveTestReq{}
				err := c.SendPkt(p, c.seqId())
				if err != nil {
					c.logger.Errorf("cmpp.active send error: %v", err)
				} else {
					atomic.AddInt32(&c.counter, 1)
				}
			}
		}
	}()
}
func (l *cmppListener) Close() error {
	return l.Listener.Close()
}
