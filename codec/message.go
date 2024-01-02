package codec

import (
	"sync/atomic"

	"github.com/zhiyin2021/zysms/smserror"
)

// 处理GSM协议TP-DCS数据编码方案
/*
 * 数据编码方案 TP-DCS（TP-Data-Coding-Scheme）长度1byte
 * Bit No.7 与Bit No.6 :一般设置为00； Bit No.5： 0—文本未压缩， 1—文本用GSM 标准压缩算法压缩
 * Bit No.4：0—表示 Bit No.1、Bit No.0 为保留位，不含信息类型信息， 1—表示Bit No.1、Bit No.0 含有信息类型信息
 * Bit No.3 与Bit No.2： 00—默认的字母表(7bit 编码) 01—8bit， 10—USC2（16bit）编码 11—预留
 * Bit No.1 与Bit No.0： 00—Class 0， 01—Class 1， 10—Class 2（SIM 卡特定信息），11—Class 3//写卡
 */
var (
	ref = uint32(0)
)

const (
	// GSM specific, short message must be no larger than 140 octets
	SM_GSM_MSG_LEN = 140
)

// type msgUDH struct {
// 	Ref   byte
// 	Total byte
// 	Seq   byte
// }

// ShortMessage message.
type ShortMessage struct {
	messageLen  byte
	message     string
	enc         Encoding
	udHeader    UDH
	messageData []byte
}

// NewLongMessage returns long message splitted into multiple short message
func NewLongMessage(message string) (s []*ShortMessage, err error) {
	enc := ASCII
	if HasWidthChar(message) {
		enc = UCS2
	}
	return NewLongMessageWithEncoding(message, enc)
}

// NewLongMessageWithEncoding returns long message splitted into multiple short message with encoding of choice
func NewLongMessageWithEncoding(message string, enc Encoding) (s []*ShortMessage, err error) {
	sm := &ShortMessage{
		message: message,
		enc:     enc,
	}
	return sm.split()
}

func (c *ShortMessage) UDHeader() UDH {
	return c.udHeader
}

// GetMessageData returns underlying binary message.
func (c *ShortMessage) GetMessageData() (d []byte) {
	return c.messageData
}
func (c *ShortMessage) IsLongMessage() bool {
	return len(c.messageData) > 3 && (c.messageData[0] == 0x05 || c.messageData[0] == 0x06) && c.messageData[1] == 0x00 && c.messageData[2] == 0x03
}

// GetMessageWithEncoding returns (decoded) underlying message.
func (c *ShortMessage) GetMessage() string {
	if c.message != "" {
		return c.message
	}
	if c.enc == nil {
		if HasWidthChar(c.message) {
			c.enc = UCS2
		} else {
			c.enc = ASCII
		}
	}
	if len(c.messageData) > 0 {
		st, _ := c.enc.Decode(c.messageData)
		return st
	}
	return ""
}
func (c *ShortMessage) GetConcatInfo() (totalParts, partNum, mref byte, found bool) {
	if c.udHeader != nil {
		return c.udHeader.GetConcatInfo()
	}
	return 0, 0, 0, false
}

// SetMessageWithEncoding sets message with encoding.
func (c *ShortMessage) SetMessage(message string, enc Encoding) (err error) {
	c.enc = enc
	if c.enc == nil {
		if HasWidthChar(message) {
			c.enc = UCS2
		} else {
			c.enc = ASCII
		}
	}
	if c.messageData, err = c.enc.Encode(message); err == nil {
		c.message = message
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

// The encoding interface can implement the Splitter interface for ad-hoc splitting rule
func (c *ShortMessage) split() (multiSM []*ShortMessage, err error) {
	if c.enc == nil {
		if HasWidthChar(c.message) {
			c.enc = UCS2
		} else {
			c.enc = ASCII
		}
	}
	gsmLen := SM_GSM_MSG_LEN
	udhLen := 6
	if c.enc == ASCII {
		gsmLen = 160
		udhLen = 7
	}
	// check if encoding implements Splitter
	// splitter, ok := c.enc.(Splitter)
	// check if encoding implements Splitter or split is necessary
	// if !ok || !splitter.ShouldSplit(c.message, gsmLen) {
	// 	err = c.SetMessage(c.message, c.enc)
	// 	multiSM = []*ShortMessage{c}
	// 	return
	// }

	// reserve 6 bytes for concat message UDH
	segments, err := c.enc.EncodeSplit(c.message, gsmLen-udhLen)
	if err != nil {
		return nil, err
	}
	// prealloc result
	multiSM = make([]*ShortMessage, 0, len(segments))
	// all segments will have the same ref id
	ref := getRefNum()

	// construct SM(s)
	for i, seg := range segments {
		// create new SM, encode data
		multiSM = append(multiSM, &ShortMessage{
			enc: c.enc,
			// message: we don't really care
			messageData: seg,
			udHeader:    UDH{NewIEConcatMessage(uint8(len(segments)), uint8(i+1), uint8(ref))}, //&msgUDH{ref, total, byte(seq + 1)}, //    UDH{NewIEConcatMessage(uint8(len(segments)), uint8(i+1), uint8(ref))},
		})
	}
	return
}

// Marshal implements PDU interface.
func (c *ShortMessage) Marshal(b *BytesWriter) {
	var (
		udhBin []byte
	)

	c.messageLen = byte(len(c.messageData))
	// Prepend UDH to message data if there are any
	if c.udHeader != nil && c.udHeader.UDHL() > 0 {
		udhBin, _ = c.udHeader.MarshalBinary()
		c.messageLen += byte(len(udhBin))
	}

	_ = b.WriteByte(byte(c.MsgLength()))
	if udhBin != nil {
		_, _ = b.Write(udhBin)
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
		udh := UDH{}
		_, err = udh.UnmarshalBinary(c.messageData)
		if err != nil {
			return
		}
		c.udHeader = udh

		f := c.udHeader.UDHL()
		if f > len(c.messageData) {
			err = smserror.ErrUDHTooLong
			return
		}

		c.messageData = c.messageData[f:]
	}

	c.enc = GetCodec(enc)
	if c.enc == nil {
		if HasWidthChar(c.message) {
			c.enc = UCS2
		} else {
			c.enc = ASCII
		}
	}
	c.message, _ = c.enc.Decode(c.messageData)
	return
}

// Encoding returns message encoding.
func (c *ShortMessage) Encoding() byte {
	if c.enc == ASCII {
		// 国内 CMPP,SMGP 不支持 gms7bit,  ascii一样以160字符计费
		return GSM7BITCoding
	}
	if c.enc == nil {
		c.enc = UCS2
	}
	return c.enc.DataCoding()
}

// returns an atomically incrementing number each time it's called
func getRefNum() uint32 {
	return atomic.AddUint32(&ref, 1)
}
