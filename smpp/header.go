package smpp

import (
	"encoding/binary"
	"sync/atomic"

	"github.com/zhiyin2021/zysms/codec"
)

func nextSequenceNumber(s *uint32) (v uint32) {
	// & 0x7FFFFFFF: cater for integer overflow
	// Allowed range is 0x01 to 0x7FFFFFFF. This
	// will still result in a single invalid value
	// of 0x00 every ~2 billion PDUs (not too bad):
	if v = atomic.AddUint32(s, 1) & 0x7FFFFFFF; v <= 0 {
		v = 1
	}
	return
}

// Header represents PDU header.
type Header struct {
	CommandLength  uint32
	CommandID      CommandId
	CommandStatus  CommandStatus
	SequenceNumber uint32
}

// ParseHeader parses PDU header.
func ParseHeader(v [16]byte) (h Header) {
	h.CommandLength = binary.BigEndian.Uint32(v[:])
	h.CommandID = CommandId(binary.BigEndian.Uint32(v[4:]))
	h.CommandStatus = CommandStatus(binary.BigEndian.Uint32(v[8:]))
	h.SequenceNumber = binary.BigEndian.Uint32(v[12:])
	return
}

// Unmarshal from buffer.
func (c *Header) Unmarshal(b *codec.BytesReader) (err error) {
	c.CommandLength = b.ReadU32()
	c.CommandID = CommandId(b.ReadU32())
	c.CommandStatus = CommandStatus(b.ReadU32())
	c.SequenceNumber = b.ReadU32()
	return b.Err()
}

var sequenceNumber uint32

// AssignSequenceNumber assigns sequence number auto-incrementally.
func (c *Header) AssignSequenceNumber() {
	c.SetSequenceNumber(nextSequenceNumber(&sequenceNumber))
}

// ResetSequenceNumber resets sequence number.
func (c *Header) ResetSequenceNumber() {
	c.SequenceNumber = 1
}

// GetSequenceNumber returns assigned sequence number.
func (c *Header) GetSequenceNumber() uint32 {
	return c.SequenceNumber
}

// SetSequenceNumber manually sets sequence number.
func (c *Header) SetSequenceNumber(v uint32) {
	c.SequenceNumber = v
}

// Marshal to buffer.
func (c *Header) Marshal(b *codec.BytesWriter) {
	b.Grow(16)
	b.WriteU32(c.CommandLength)
	b.WriteU32(uint32(c.CommandID))
	b.WriteU32(uint32(c.CommandStatus))
	b.WriteU32(c.SequenceNumber)
}
