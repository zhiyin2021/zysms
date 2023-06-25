package proto

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/zhiyin2021/zysms/utils"
)

// Common errors.
var ErrMethodParamsInvalid = errors.New("params passed to method is invalid")

// Protocol errors.
var ErrTotalLengthInvalid = errors.New("total_length in Packet data is invalid")
var ErrCommandIdInvalid = errors.New("command_Id in Packet data is invalid")
var ErrCommandIdNotSupported = errors.New("command_Id in Packet data is not supported")

type Packer interface {
	Pack(seqId uint32) []byte
	Unpack(data []byte)
}
type packet struct {
	data  []byte
	index int
}

// packet必须按照顺序读取或写入，否则会出错
func NewPacket(data []byte) *packet {
	return &packet{data: data, index: 0}
}

func (c *packet) ReadBytes(count int) []byte {
	defer func() {
		c.index += count
	}()
	last := c.index + count
	return c.data[c.index:last]
}
func (c *packet) ReadByte() byte {
	defer func() {
		c.index++
	}()
	return c.data[c.index]
}
func (c *packet) ReadU32() uint32 {
	defer func() {
		c.index += 4
	}()
	return binary.BigEndian.Uint32(c.data[c.index : c.index+4])
}

func (c *packet) ReadU64() uint64 {
	defer func() {
		c.index += 8
	}()
	return binary.BigEndian.Uint64(c.data[c.index : c.index+8])
}

func (c *packet) ReadStr(count int) string {
	defer func() {
		c.index += count
	}()
	last := c.index + count
	return string(bytes.TrimLeft(c.data[c.index:last], "\x00"))
}

func (c *packet) Skip(count int) {
	c.index += count
}

func (c *packet) WriteBytes(buf []byte, count int) {
	last := c.index + count
	defer func() {
		c.index = last
	}()
	copy(c.data[c.index:last], buf)
}

func (c *packet) WriteByte(b byte) {
	defer func() {
		c.index++
	}()
	c.data[c.index] = b
}
func (c *packet) WriteU32(n uint32) {
	defer func() {
		c.index += 4
	}()
	binary.BigEndian.PutUint32(c.data[c.index:c.index+4], n)
}

func (c *packet) WriteU64(n uint64) {
	defer func() {
		c.index += 8
	}()
	binary.BigEndian.PutUint64(c.data[c.index:c.index+8], n)
}
func (c *packet) WriteStr(data string, count int) {
	last := c.index + count
	defer func() {
		c.index = last
	}()
	tmp := utils.OctetString(data, count)
	copy(c.data[c.index:last], tmp)
}
func (c *packet) Index() int {
	return c.index
}