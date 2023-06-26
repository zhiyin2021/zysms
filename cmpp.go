package zysms

import (
	"encoding/binary"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms/cmpp"
	"github.com/zhiyin2021/zysms/enum"
	"github.com/zhiyin2021/zysms/proto"
)

type cmppConn struct {
	net.Conn
	State enum.State
	Typ   cmpp.Version
	// for SeqId generator goroutine
	// SeqId  <-chan uint32
	// done   chan<- struct{}
	_seqId  uint32
	done    chan struct{}
	counter int32
	logger  *logrus.Entry
}

// New returns an abstract structure for successfully
// established underlying net.Conn.
func newCmppConn(conn net.Conn, typ cmpp.Version) *cmppConn {
	c := &cmppConn{
		Conn:   conn,
		Typ:    typ,
		_seqId: 0,
		done:   make(chan struct{}, 1),
		logger: logrus.WithFields(logrus.Fields{"r": conn.RemoteAddr(), "v": typ}),
	}
	tc := c.Conn.(*net.TCPConn)
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(1 * time.Minute) // 1min

	return c
}
func (c *cmppConn) Logger() *logrus.Entry {
	return c.logger
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
	if c.Typ == cmpp.V30 {
		if rsp, ok := p.(*cmpp.Cmpp3ConnRsp); ok {
			status = uint8(rsp.Status)
		} else {
			return enum.ErrRespNotMatch
		}
	} else {
		if rsp, ok := p.(*cmpp.Cmpp2ConnRsp); ok {
			status = rsp.Status
		} else {
			return enum.ErrRespNotMatch
		}
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
		// close(c.done)  // let the SeqId goroutine exit.
		c.Conn.Close() // close the underlying net.Conn
		c.State = enum.CONN_CLOSED
		c.logger.Infoln("cmpp.conn close")
		c.done <- struct{}{}
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
		return enum.ErrConnIsClosed
	}

	if seqId == 0 {
		seqId = c.seqId()
	}
	data := pkt.Pack(seqId)

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

// RecvAndUnpackPkt receives cmpp byte stream, and unpack it to some cmpp packet structure.
func (c *cmppConn) RecvPkt(timeout time.Duration) (proto.Packer, error) {
	if c.State == enum.CONN_CLOSED {
		return nil, enum.ErrConnIsClosed
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
		c.logger.Errorf("cmpp.recv read totalLen error: %v ", err)
		return nil, err
	}

	if c.Typ == cmpp.V30 && (rb.totalLen < cmpp.CMPP3_PACKET_MIN || rb.totalLen > cmpp.CMPP3_PACKET_MAX) {
		return nil, proto.ErrTotalLengthInvalid
	} else if rb.totalLen < cmpp.CMPP2_PACKET_MIN || rb.totalLen > cmpp.CMPP2_PACKET_MAX {
		return nil, proto.ErrTotalLengthInvalid
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
				return nil, enum.ErrReadCmdIDTimeout
			}
		}
		return nil, err
	}

	if !((rb.commandId > cmpp.CMPP_REQUEST_MIN && rb.commandId < cmpp.CMPP_REQUEST_MAX) ||
		(rb.commandId > cmpp.CMPP_RESPONSE_MIN && rb.commandId < cmpp.CMPP_RESPONSE_MAX)) {
		return nil, proto.ErrCommandIdInvalid
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
		c.logger.Infof("%d => %s", p.SeqId(), rb.commandId.String())
		if rb.commandId == cmpp.CMPP_ACTIVE_TEST {
			resp := &cmpp.CmppActiveTestRsp{}
			c.SendPkt(resp, p.SeqId())
			atomic.AddInt32(&c.counter, -1)
		}
		return p, nil
	}
	return nil, proto.ErrCommandIdNotSupported

}

type CmppListener struct {
	net.Listener
	typ cmpp.Version
}

func newCmppListener(l net.Listener, v cmpp.Version) *CmppListener {
	return &CmppListener{l, v}
}

func (l *CmppListener) Accept() (SmsConn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	conn := newCmppConn(c, l.typ)
	conn.SetState(enum.CONN_CONNECTED)
	conn.startActiveTest()
	return conn, nil
}
func (c *cmppConn) startActiveTest() {
	go func() {
		t := time.NewTicker(10 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-c.done:
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
func (l *CmppListener) Close() error {
	return l.Listener.Close()
}
