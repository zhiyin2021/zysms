package cmpp

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zhiyin2021/zysms/cmpp/codec"
	"github.com/zhiyin2021/zysms/proto"
)

type State uint8

// Errors for conn operations
var (
	ErrConnIsClosed       = errors.New("connection is closed")
	ErrReadCmdIDTimeout   = errors.New("read commandId timeout")
	ErrReadPktBodyTimeout = errors.New("read packet body timeout")
)

var noDeadline = time.Time{}

// Conn States
const (
	CONN_CLOSED State = iota
	CONN_CONNECTED
	CONN_AUTHOK
)

type Conn struct {
	net.Conn
	State State
	Typ   codec.Version
	// for SeqId generator goroutine
	// SeqId  <-chan uint32
	// done   chan<- struct{}
	_seqId uint32
}

// New returns an abstract structure for successfully
// established underlying net.Conn.
func newConn(conn net.Conn, typ codec.Version) *Conn {
	// seqId, done := newSeqIdGenerator()
	c := &Conn{
		Conn: conn,
		Typ:  typ,
		// SeqId:  seqId,
		// done:   done,
		_seqId: 0,
	}
	tc := c.Conn.(*net.TCPConn) // Always tcpconn
	tc.SetKeepAlive(true)       //Keepalive as default
	return c
}

func (c *Conn) Close() {
	if c != nil {
		if c.State == CONN_CLOSED {
			return
		}
		// close(c.done)  // let the SeqId goroutine exit.
		c.Conn.Close() // close the underlying net.Conn
		c.State = CONN_CLOSED
	}
}

func (c *Conn) SetState(state State) {
	c.State = state
}
func (c *Conn) SeqId() uint32 {
	return atomic.AddUint32(&c._seqId, 1)
}

// SendPkt pack the cmpp packet structure and send it to the other peer.
func (c *Conn) SendPkt(packet proto.Packer, seqId uint32) error {
	if c.State == CONN_CLOSED {
		return ErrConnIsClosed
	}

	data := packet.Pack(seqId)

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
	commandId codec.CommandId
	leftData  [defaultReadBufferSize]byte
}

var readBufferPool = sync.Pool{
	New: func() interface{} {
		return &readBuffer{}
	},
}

// RecvAndUnpackPkt receives cmpp byte stream, and unpack it to some cmpp packet structure.
func (c *Conn) RecvAndUnpackPkt(timeout time.Duration) (interface{}, error) {
	if c.State == CONN_CLOSED {
		return nil, ErrConnIsClosed
	}
	defer c.SetReadDeadline(noDeadline)

	rb := readBufferPool.Get().(*readBuffer)
	defer readBufferPool.Put(rb)

	// Total_Length in packet
	if timeout != 0 {
		c.SetReadDeadline(time.Now().Add(timeout))
	}
	err := binary.Read(c.Conn, binary.BigEndian, &rb.totalLen)
	if err != nil {
		return nil, err
	}

	if c.Typ == codec.V30 {
		if rb.totalLen < codec.CMPP3_PACKET_MIN || rb.totalLen > codec.CMPP3_PACKET_MAX {
			return nil, proto.ErrTotalLengthInvalid
		}
	}

	if c.Typ == codec.V21 || c.Typ == codec.V20 {
		if rb.totalLen < codec.CMPP2_PACKET_MIN || rb.totalLen > codec.CMPP2_PACKET_MAX {
			return nil, proto.ErrTotalLengthInvalid
		}
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
				return nil, ErrReadCmdIDTimeout
			}
		}
		return nil, err
	}

	if !((rb.commandId > codec.CMPP_REQUEST_MIN && rb.commandId < codec.CMPP_REQUEST_MAX) ||
		(rb.commandId > codec.CMPP_RESPONSE_MIN && rb.commandId < codec.CMPP_RESPONSE_MAX)) {
		return nil, proto.ErrCommandIdInvalid
	}

	// The left packet data (start from seqId in header).
	if timeout != 0 {
		c.SetReadDeadline(time.Now().Add(timeout))
	}
	var leftData = rb.leftData[0:(rb.totalLen - 8)]
	_, err = io.ReadFull(c.Conn, leftData)
	if err != nil {
		netErr, ok := err.(net.Error)
		if ok {
			if netErr.Timeout() {
				return nil, ErrReadPktBodyTimeout
			}
		}
		return nil, err
	}

	var p proto.Packer
	switch rb.commandId {
	case codec.CMPP_CONNECT:
		p = &codec.CmppConnReq{}
	case codec.CMPP_CONNECT_RESP:
		if c.Typ == codec.V30 {
			p = &codec.Cmpp3ConnRsp{}
		} else {
			p = &codec.Cmpp2ConnRsp{}
		}
	case codec.CMPP_TERMINATE:
		p = &codec.CmppTerminateReq{}
	case codec.CMPP_TERMINATE_RESP:
		p = &codec.CmppTerminateRsp{}
	case codec.CMPP_SUBMIT:
		if c.Typ == codec.V30 {
			p = &codec.Cmpp3SubmitReq{}
		} else {
			p = &codec.Cmpp2SubmitReq{}
		}
	case codec.CMPP_SUBMIT_RESP:
		if c.Typ == codec.V30 {
			p = &codec.Cmpp3SubmitRsp{}
		} else {
			p = &codec.Cmpp2SubmitRsp{}
		}
	case codec.CMPP_DELIVER:
		if c.Typ == codec.V30 {
			p = &codec.Cmpp3DeliverReq{}
		} else {
			p = &codec.Cmpp2DeliverReq{}
		}
	case codec.CMPP_DELIVER_RESP:
		if c.Typ == codec.V30 {
			p = &codec.Cmpp3DeliverRsp{}
		} else {
			p = &codec.Cmpp2DeliverRsp{}
		}
	case codec.CMPP_FWD:
		if c.Typ == codec.V30 {
			p = &codec.Cmpp3FwdReq{}
		} else {
			p = &codec.Cmpp2FwdReq{}
		}
	case codec.CMPP_FWD_RESP:
		if c.Typ == codec.V30 {
			p = &codec.Cmpp3FwdRsp{}
		} else {
			p = &codec.Cmpp2FwdRsp{}
		}
	case codec.CMPP_ACTIVE_TEST:
		p = &codec.CmppActiveTestReq{}
	case codec.CMPP_ACTIVE_TEST_RESP:
		p = &codec.CmppActiveTestRsp{}

	default:
		p = nil
		return nil, proto.ErrCommandIdNotSupported
	}

	p.Unpack(leftData)
	return p, nil
}
