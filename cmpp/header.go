package cmpp

import (
	"encoding/binary"
	"fmt"
	"sync/atomic"

	"github.com/zhiyin2021/zysms/codec"
)

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
	SequenceNumber int32
}

func (c Header) String() string {
	return fmt.Sprintf("{len:%d,cmd:%#v,seq:%d}", c.CommandLength, c.CommandID, c.SequenceNumber)
}
func (c *Header) GetCommandID() uint32 {
	return uint32(c.CommandID)
}

func (c *Header) GetCommandLength() uint32 {
	return c.CommandLength
}

// ParseHeader parses PDU header.
func ParseHeader(v [12]byte) (h Header) {
	h.CommandLength = binary.BigEndian.Uint32(v[:])
	h.CommandID = codec.CommandId(binary.BigEndian.Uint32(v[4:]))
	h.SequenceNumber = int32(binary.BigEndian.Uint32(v[8:]))
	return
}

// Unmarshal from buffer.
func (c *Header) Unmarshal(b *codec.BytesReader) (err error) {
	c.CommandLength = b.ReadU32()
	c.CommandID = codec.CommandId(b.ReadU32())
	c.SequenceNumber = int32(b.ReadU32())
	return b.Err()
}

var sequenceNumber int32

// AssignSequenceNumber assigns sequence number auto-incrementally.
func (c *Header) AssignSequenceNumber() {
	c.SetSequenceNumber(nextSequenceNumber(&sequenceNumber))
}

// ResetSequenceNumber resets sequence number.
func (c *Header) ResetSequenceNumber() {
	c.SequenceNumber = 1
}

// GetSequenceNumber returns assigned sequence number.
func (c *Header) GetSequenceNumber() int32 {
	return c.SequenceNumber
}

// SetSequenceNumber manually sets sequence number.
func (c *Header) SetSequenceNumber(v int32) {
	c.SequenceNumber = v
}

// Marshal to buffer.
func (c *Header) Marshal(b *codec.BytesWriter) {
	b.Grow(12)
	b.WriteU32(c.CommandLength)
	b.WriteU32(uint32(c.CommandID))
	b.WriteU32(uint32(c.SequenceNumber))
}
