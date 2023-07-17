package codec

import (
	"bytes"
	"fmt"
)

// ByteBuffer wraps over bytes.Buffer with additional features.
type ByteBuffer struct {
	*bytes.Buffer
}

// NewBuffer create new buffer from preallocated buffer array.
func NewBuffer(inp []byte) *ByteBuffer {
	if inp == nil {
		inp = make([]byte, 0, 512)
	}
	return &ByteBuffer{Buffer: bytes.NewBuffer(inp)}
}

// ReadN read n-bytes from buffer.
func (c *ByteBuffer) ReadN(n int) (r []byte, err error) {
	if n > 0 {
		if c.Len() >= n { // optimistic branching
			r = make([]byte, n)
			_, _ = c.Read(r)
		} else {
			err = ErrBufferNotEnoughByteToRead
		}
	}
	return
}

// ReadShort reads short from buffer.
func (c *ByteBuffer) ReadShort() (r int16, err error) {
	v, err := c.ReadN(SizeShort)
	if err == nil {
		r = int16(endianese.Uint16(v))
	}
	return
}

// WriteShort writes short to buffer.
func (c *ByteBuffer) WriteShort(v int16) {
	var b [SizeShort]byte
	endianese.PutUint16(b[:], uint16(v))
	_, _ = c.Write(b[:])
}

// ReadInt reads int from buffer.
func (c *ByteBuffer) ReadInt() (r int32, err error) {
	v, err := c.ReadN(SizeInt)
	if err == nil {
		r = int32(endianese.Uint32(v))
	}
	return
}

// WriteInt writes int to buffer.
func (c *ByteBuffer) WriteInt(v int32) {
	var b [SizeInt]byte
	endianese.PutUint32(b[:], uint32(v))
	_, _ = c.Write(b[:])
}

// WriteBuffer appends buffer.
func (c *ByteBuffer) WriteBuffer(d *ByteBuffer) {
	if d != nil {
		_, _ = c.Write(d.Bytes())
	}
}

func (c *ByteBuffer) writeString(st string, isCString bool, enc Encoding) (err error) {
	if len(st) > 0 {
		var payload []byte
		if payload, err = enc.Encode(st); err == nil {
			_, _ = c.Write(payload)
		}
	}

	if err == nil && isCString {
		_ = c.WriteByte(0)
	}

	return
}

// WriteCString writes c-string.
func (c *ByteBuffer) WriteCString(s string) error {
	return c.writeString(s, true, ASCII)
}

// WriteCStringWithEnc write c-string with encoding.
func (c *ByteBuffer) WriteCStringWithEnc(s string, enc Encoding) error {
	return c.writeString(s, true, enc)
}

// ReadCString read c-string.
func (c *ByteBuffer) ReadCString() (st string, err error) {
	buf, err := c.ReadBytes(0)
	if err == nil && len(buf) > 0 { // optimistic branching
		st = string(buf[:len(buf)-1])
	}
	return
}

// HexDump returns hex dump.
func (c *ByteBuffer) HexDump() string {
	return fmt.Sprintf("%x", c.Buffer.Bytes())
}
