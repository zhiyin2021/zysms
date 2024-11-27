package codec

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
	endianese                    = binary.BigEndian
)

type IBuffer interface {
	Reset()
	Write([]byte) (int, error)
}

// bytesBuffer wraps over bytes.Buffer with additional features.
type bytesBuffer struct {
	*bytes.Buffer
	err error
}
type BytesReader struct {
	bytes []byte
	*bytesBuffer
}
type BytesWriter struct {
	*bytesBuffer
}

// type SyncPool[T IBuffer] struct {
// 	Total int
// 	pool  sync.Pool
// }

// func (sp *SyncPool[T]) Get(buf ...[]byte) T {
// 	obj := sp.pool.Get().(T)
// 	obj.Reset()
// 	if len(buf) > 0 {
// 		obj.Write(buf[0])
// 	}
// 	return obj
// }

// func (sp *SyncPool[T]) Put(obj T) {
// 	sp.pool.Put(obj)
// }

// var (
// 	WriterPool = NewWirterPool()
// 	ReaderPool = NewReaderPool()
// )

// func NewWirterPool() *SyncPool[*BytesWriter] {
// 	sp := &SyncPool[*BytesWriter]{}
// 	sp.pool.New = func() any {
// 		sp.Total++
// 		return &BytesWriter{bytesBuffer: &bytesBuffer{Buffer: bytes.NewBuffer(make([]byte, 0, 12)), err: nil}}
// 	}
// 	return sp
// }

// func NewReaderPool() *SyncPool[*BytesReader] {
// 	sp := &SyncPool[*BytesReader]{}
// 	sp.pool.New = func() any {
// 		sp.Total++
// 		return &BytesReader{bytesBuffer: &bytesBuffer{Buffer: bytes.NewBuffer(make([]byte, 0, 12)), err: nil}}
// 	}
// 	return sp
// }

func NewWriter() *BytesWriter {
	return &BytesWriter{bytesBuffer: &bytesBuffer{Buffer: bytes.NewBuffer([]byte{}), err: nil}}

}

func NewReader(buf []byte) *BytesReader {
	return &BytesReader{bytesBuffer: &bytesBuffer{Buffer: bytes.NewBuffer(buf), err: nil}}
}

// ReadN read n-bytes from buffer.
func (c *BytesReader) ReadN(n int) (r []byte) {
	if c.err == nil {
		if n > 0 {
			if c.Len() >= n { // optimistic branching
				r = make([]byte, n)
				_, _ = c.Read(r)
			} else {
				c.err = fmt.Errorf("not enough byte to read from buffer(%d>%d): %x", n, c.Len(), c.bytes)
				panic(c.err)
			}
		}
	}
	return
}

// ReadShort reads short from buffer.
func (c *BytesReader) ReadU16() (r uint16) {
	if c.err == nil {
		v := c.ReadN(SizeShort)
		if c.err == nil {
			r = endianese.Uint16(v)
		}
	}
	return
}

// WriteShort writes short to buffer.
func (c *BytesWriter) WriteU16(v uint16) {
	var b [SizeShort]byte
	endianese.PutUint16(b[:], v)
	_, _ = c.Write(b[:])
}
func (c *BytesReader) ReadU8() byte {
	if c.err == nil {
		var v byte
		v, c.err = c.Buffer.ReadByte()
		return v
	}
	return 0
}

// ReadInt reads int from buffer.
func (c *BytesReader) ReadU32() (r uint32) {
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
func (c *BytesReader) ReadU64() (r uint64) {
	if c.err == nil {
		v := c.ReadN(SizeLong)
		if c.err == nil {
			r = endianese.Uint64(v)
		}
	}
	return
}

// WriteInt writes int to buffer.
func (c *BytesWriter) WriteU32(v uint32) {
	var b [SizeInt]byte
	endianese.PutUint32(b[:], v)
	_, _ = c.Write(b[:])
}

// WriteInt writes int to buffer.
func (c *BytesWriter) WriteU64(v uint64) {
	var b [SizeLong]byte
	endianese.PutUint64(b[:], v)
	_, _ = c.Write(b[:])
}

func (c *BytesWriter) writeString(st string, isCString bool, enc Encoding, count int) (err error) {
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

// WriteStr
func (c *BytesWriter) WriteStr(s string, count int) error {
	payload := []byte(s)
	if len(payload) > count {
		payload = payload[:count]
	} else {
		payload = append(payload, make([]byte, count-len(payload))...)
	}
	_, c.err = c.Write(payload)
	return c.err
}

// WriteCString writes c-string.
func (c *BytesWriter) WriteCStr(s string) error {
	return c.writeString(s, true, LATIN1, 0)
}

// WriteCStringWithEnc write c-string with encoding.
func (c *BytesWriter) WriteCStrWithEnc(s string, enc Encoding) error {
	return c.writeString(s, true, enc, 0)
}

// ReadStr
func (c *BytesReader) ReadStr(count int) string {
	buf := c.ReadN(count)
	if c.err == nil && len(buf) > 0 { // optimistic branching
		return string(bytes.TrimRight(buf, "\x00"))
	}
	return ""
}

// WriteInt writes int to buffer.
func (c *BytesReader) WriteBytes(buf []byte) {
	c.bytes = append([]byte{}, buf...)
	c.Write(buf)
}

// ReadCString
func (c *BytesReader) ReadCStr() (st string) {
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
