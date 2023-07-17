package codec

import (
	"log"
	"sync/atomic"
)

var (
	ref = uint32(0)
)

type msgUDH struct {
	Ref   byte
	Total byte
	Seq   byte
}

// ShortMessage message.
type ShortMessage struct {
	messageLen  byte
	message     string
	enc         Encoding
	udHeader    *msgUDH
	messageData []byte
}

func (c *ShortMessage) UDHeader() *msgUDH {
	return c.udHeader
}

// GetMessageData returns underlying binary message.
func (c *ShortMessage) GetMessageData() (d []byte) {
	return c.messageData
}

// GetMessageWithEncoding returns (decoded) underlying message.
func (c *ShortMessage) GetMessage() string {
	if c.message != "" {
		return c.message
	}
	if c.enc == nil {
		c.enc = ASCII
	}
	if len(c.messageData) > 0 {
		st, _ := c.enc.Decode(c.messageData)
		return st
	}
	return ""
}

// SetMessageWithEncoding sets message with encoding.
func (c *ShortMessage) SetMessage(message string, enc Encoding) (err error) {
	if c.messageData, err = enc.Encode(message); err == nil {
		c.message = message
		c.enc = enc
	}
	return
}

func (c *ShortMessage) MsgLength() int {
	n := len(c.messageData)
	if c.udHeader != nil {
		n += 6
	}
	return n
}

// The encoding interface can implement the data.Splitter interface for ad-hoc splitting rule
func (c *ShortMessage) Split() (multiSM []*ShortMessage, err error) {
	if c.enc == nil {
		c.enc = ASCII
	}
	maxLen := uint(140)
	if c.enc == ASCII {
		maxLen = 160
	}
	// check if encoding implements data.Splitter
	splitter, ok := c.enc.(Splitter)
	// check if encoding implements data.Splitter or split is necessary
	if !ok || !splitter.ShouldSplit(c.message, maxLen) {
		multiSM = []*ShortMessage{c}
		return
	}

	// reserve 6 bytes for concat message UDH
	segments, err := splitter.EncodeSplit(c.message, maxLen-6)
	if err != nil {
		return nil, err
	}
	// prealloc result
	multiSM = make([]*ShortMessage, 0, len(segments))
	// all segments will have the same ref id
	ref := byte(getRefNum())
	total := byte(len(segments))
	// construct SM(s)
	for seq, seg := range segments {
		// create new SM, encode data
		multiSM = append(multiSM, &ShortMessage{
			enc: c.enc,
			// message: we don't really care
			messageData: seg,
			udHeader:    &msgUDH{ref, total, byte(seq)}, //    UDH{NewIEConcatMessage(uint8(len(segments)), uint8(i+1), uint8(ref))},
		})
		log.Printf("%d => %d", seq, len(seg))
	}
	return
}

// Marshal implements PDU interface.
func (c *ShortMessage) Marshal(b *BytesWriter) {
	c.messageLen = byte(len(c.messageData))
	_ = b.WriteByte(byte(c.MsgLength()))
	if c.udHeader != nil {
		buf := []byte{0x05, 0x00, 0x03, c.udHeader.Ref, c.udHeader.Total, c.udHeader.Seq}
		_, _ = b.Write(buf)
	}
	// short_message
	_, _ = b.Write(c.messageData)
}

// Unmarshal implements PDU interface.
func (c *ShortMessage) Unmarshal(b *BytesReader, udhi bool, enc byte) (err error) {
	c.messageLen = b.ReadByte()
	c.messageData = b.ReadN(int(c.messageLen))
	// If short message length is non zero, short message contains User-Data Header
	// Else UDH should be in TLV field MessagePayload
	if udhi && c.messageLen > 0 {
		n := c.messageData[0] + 1
		if n > c.messageLen {
			if n == 6 {
				c.udHeader = &msgUDH{c.messageData[4], c.messageData[5], c.messageData[6]}
				// 0x05 数据头总长度
				// 0x00 信息标识
				// 0x04 头信息长度
				// 0x00 信息序列号
				// 0x00 总条数
				// 0x01 当前条数
				c.messageData = c.messageData[n:]
			}
		}
	}
	c.enc = GetCodec(enc)
	if c.enc == nil {
		c.enc = UCS2
	}
	c.message, _ = c.enc.Decode(c.messageData)
	return
}

// Encoding returns message encoding.
func (c *ShortMessage) Encoding() Encoding {
	return c.enc
}

// returns an atomically incrementing number each time it's called
func getRefNum() uint32 {
	return atomic.AddUint32(&ref, 1)
}
