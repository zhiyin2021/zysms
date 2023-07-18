package sgip

import (
	"encoding/binary"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/zhiyin2021/zysms/codec"
)

type SgipHeader struct {
	SeqId [3]uint32 // 源节点编号 + 月日时分秒 + 流水序号
}

func getTm() uint32 {
	tm, _ := strconv.ParseUint(time.Now().Format("0215040506"), 10, 10)
	return uint32(tm)
}

func nextSequenceNumber(s *int32) (v int32) {
	// & 0x7FFFFFFF: cater for integer overflow
	// Allowed range is 0x01 to 0x7FFFFFFF. This
	// will still result in a single invalid value
	// of 0x00 every ~2 billion PDUs (not too bad):
	if v = atomic.AddInt32(s, 1) & 0x7FFFFFFF; v <= 0 {
		v = 1
	}
	return
}

// Header represents PDU header.
type Header struct {
	CommandLength  uint32
	CommandID      codec.CommandId
	SequenceNumber [3]uint32
	NodeId         uint32
}

// ParseHeader parses PDU header.
func ParseHeader(v [20]byte) (h Header) {
	h.CommandLength = binary.BigEndian.Uint32(v[:])
	h.CommandID = codec.CommandId(binary.BigEndian.Uint32(v[4:]))
	h.SequenceNumber[0] = binary.BigEndian.Uint32(v[8:])
	h.SequenceNumber[1] = binary.BigEndian.Uint32(v[12:])
	h.SequenceNumber[2] = binary.BigEndian.Uint32(v[16:])
	return
}

// Unmarshal from buffer.
func (h *Header) Unmarshal(b *codec.BytesReader) error {
	h.CommandLength = b.ReadU32()
	h.CommandID = codec.CommandId(b.ReadU32())
	h.SequenceNumber[0] = b.ReadU32()
	h.SequenceNumber[1] = b.ReadU32()
	h.SequenceNumber[2] = b.ReadU32()
	return b.Err()
}

var sequenceNumber int32

// AssignSequenceNumber assigns sequence number auto-incrementally.
func (c *Header) AssignSequenceNumber() {
	c.SetSequenceNumber(nextSequenceNumber(&sequenceNumber))
}

// ResetSequenceNumber resets sequence number.
func (c *Header) ResetSequenceNumber() {
	c.SequenceNumber[2] = 1
}

// GetSequenceNumber returns assigned sequence number.
func (c *Header) GetSequenceNumber() int32 {
	return int32(c.SequenceNumber[2])
}

// SetSequenceNumber manually sets sequence number.
func (c *Header) SetSequenceNumber(v int32) {
	c.SequenceNumber[2] = uint32(v)
}

// Marshal to buffer.
func (c *Header) Marshal(b *codec.BytesWriter) {
	b.Grow(16)
	b.WriteU32(c.CommandLength)
	b.WriteU32(uint32(c.CommandID))
	b.WriteU32(c.SequenceNumber[0])
	b.WriteU32(c.SequenceNumber[1])
	b.WriteU32(c.SequenceNumber[2])
}
