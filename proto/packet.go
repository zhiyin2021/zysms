package proto

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/zhiyin2021/zysms/event"
)

const (
	// SizeByte is size of byte.
	SizeByte = 1
	// SizeShort is size of short.
	SizeShort = 2
	// SizeInt is size of int.
	SizeInt = 4
	// SizeLong is size of long.
	SizeLong = 8
)

var (
	// ErrBufferNotEnoughByteToRead indicates not enough byte(s) to read from buffer.
	ErrBufferNotEnoughByteToRead = fmt.Errorf("not enough byte to read from buffer")

	endianese = binary.BigEndian
)

type Packer interface {
	SeqId() uint32
	Pack(uint32, SmsProto) []byte
	Unpack([]byte, SmsProto) error
	Event() event.SmsEvent
}

// bytesBuffer wraps over bytes.Buffer with additional features.
type bytesBuffer struct {
	*bytes.Buffer
	err error
}

// NewBuffer create new buffer from preallocated buffer array.
func NewBuffer(inp []byte) *bytesBuffer {
	if inp == nil {
		inp = make([]byte, 0, 512)
	}
	return &bytesBuffer{Buffer: bytes.NewBuffer(inp), err: nil}
}
func NewCmppBuffer(dataLen, commandId, seqId uint32) *bytesBuffer {
	inp := make([]byte, 0, dataLen)
	b := &bytesBuffer{Buffer: bytes.NewBuffer(inp), err: nil}
	b.WriteU32(dataLen)
	b.WriteU32(commandId)
	b.WriteU32(seqId)
	return b
}

// ReadN read n-bytes from buffer.
func (c *bytesBuffer) ReadN(n int) (r []byte) {
	if c.err == nil {
		if n > 0 {
			if c.Len() >= n { // optimistic branching
				r = make([]byte, n)
				_, _ = c.Read(r)
			} else {
				c.err = ErrBufferNotEnoughByteToRead
			}
		}
	}
	return
}

// ReadShort reads short from buffer.
func (c *bytesBuffer) ReadU16() (r uint16) {
	if c.err == nil {
		v := c.ReadN(SizeShort)
		if c.err == nil {
			r = endianese.Uint16(v)
		}
	}
	return
}

// WriteShort writes short to buffer.
func (c *bytesBuffer) WriteU16(v uint16) {
	var b [SizeShort]byte
	endianese.PutUint16(b[:], v)
	_, _ = c.Write(b[:])
}
func (c *bytesBuffer) ReadByte() byte {
	if c.err == nil {
		var v byte
		v, c.err = c.Buffer.ReadByte()
		return v
	}
	return 0
}

// ReadInt reads int from buffer.
func (c *bytesBuffer) ReadU32() (r uint32) {
	if c.err == nil {
		v := c.ReadN(SizeInt)
		if c.err == nil {
			r = endianese.Uint32(v)
		}
	}
	return
}
func (c *bytesBuffer) Err() error {
	return c.err
}

// ReadInt reads int from buffer.
func (c *bytesBuffer) ReadU64() (r uint64) {
	if c.err == nil {
		v := c.ReadN(SizeLong)
		if c.err == nil {
			r = endianese.Uint64(v)
		}
	}
	return
}

// WriteInt writes int to buffer.
func (c *bytesBuffer) WriteU32(v uint32) {
	var b [SizeInt]byte
	endianese.PutUint32(b[:], v)
	_, _ = c.Write(b[:])
}

// WriteInt writes int to buffer.
func (c *bytesBuffer) WriteBytes(buf []byte) {
	c.Write(buf)
}

// WriteInt writes int to buffer.
func (c *bytesBuffer) WriteU64(v uint64) {
	var b [SizeLong]byte
	endianese.PutUint64(b[:], v)
	_, _ = c.Write(b[:])
}

// WriteBuffer appends buffer.
func (c *bytesBuffer) WriteBuffer(d *bytesBuffer) {
	if d != nil {
		_, _ = c.Write(d.Bytes())
	}
}

func (c *bytesBuffer) writeString(st string, isCString bool, enc Encoding, count int) (err error) {
	if len(st) > 0 {
		var payload []byte
		if payload, err = enc.Encode(st); err == nil {
			if count > 0 {
				if len(payload) > count {
					payload = payload[:count]
				} else {
					payload = append(payload, make([]byte, count-len(payload))...)
				}
			}
			_, _ = c.Write(payload)
		}
	} else {
		if count > 0 {
			_, _ = c.Write(make([]byte, count))
		}
	}

	if err == nil && isCString && count == 0 {
		_ = c.WriteByte(0)
	}
	return
}

// WriteCString writes c-string.
func (c *bytesBuffer) WriteCStrN(s string, count int) error {
	return c.writeString(s, false, ASCII, count)
}

// WriteCString writes c-string.
func (c *bytesBuffer) WriteCStr(s string) error {
	return c.writeString(s, true, ASCII, 0)
}

// WriteCStringWithEnc write c-string with encoding.
func (c *bytesBuffer) WriteCStrWithEnc(s string, enc Encoding) error {
	return c.writeString(s, true, enc, 0)
}

// ReadCString read c-string.
func (c *bytesBuffer) ReadCStrN(count int) string {
	buf := c.ReadN(count)
	if c.err == nil && len(buf) > 0 { // optimistic branching
		return string(bytes.TrimLeft(buf, "\x00"))
	}
	return ""
}

// ReadCString read c-string.
func (c *bytesBuffer) ReadCStr() (st string) {
	buf, err := c.ReadBytes(0)
	if err == nil && len(buf) > 0 { // optimistic branching
		st = string(buf[:len(buf)-1])
	}
	return
}

// HexDump returns hex dump.
func (c *bytesBuffer) HexDump() string {
	return fmt.Sprintf("%x", c.Buffer.Bytes())
}
