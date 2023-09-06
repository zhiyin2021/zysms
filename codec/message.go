package codec

import (
	"sync/atomic"
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

// NewLongMessage returns long message splitted into multiple short message
func NewLongMessage(message string) (s []*ShortMessage, err error) {
	enc := ASCII
	if hasWidthChar(message) {
		enc = UCS2
	}
	if _, err = GSM7BIT.Encode(message); err == nil {
		enc = GSM7BIT
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

func (c *ShortMessage) UDHeader() *msgUDH {
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
		if hasWidthChar(c.message) {
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

// SetMessageWithEncoding sets message with encoding.
func (c *ShortMessage) SetMessage(message string, enc Encoding) (err error) {
	c.enc = enc
	if c.enc == nil {
		if hasWidthChar(message) {
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
		if hasWidthChar(c.message) {
			c.enc = UCS2
		} else {
			c.enc = ASCII
		}
	}

	// check if encoding implements Splitter
	splitter, ok := c.enc.(Splitter)
	// check if encoding implements Splitter or split is necessary
	if !ok || !splitter.ShouldSplit(c.message, SM_GSM_MSG_LEN) {
		err = c.SetMessage(c.message, c.enc)
		multiSM = []*ShortMessage{c}
		return
	}

	// reserve 6 bytes for concat message UDH
	segments, err := splitter.EncodeSplit(c.message, SM_GSM_MSG_LEN-6)
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
			udHeader:    &msgUDH{ref, total, byte(seq + 1)}, //    UDH{NewIEConcatMessage(uint8(len(segments)), uint8(i+1), uint8(ref))},
		})
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
		if n == 6 && n < c.messageLen {
			c.udHeader = &msgUDH{
				c.messageData[3],
				c.messageData[4],
				c.messageData[5],
			}
			// 0x05 数据头总长度
			// 0x00 信息标识
			// 0x04 头信息长度
			// 0x00 信息序列号
			// 0x00 总条数
			// 0x01 当前条数
			c.messageData = c.messageData[n:]
		}
	}
	c.enc = GetCodec(enc)
	if c.enc == nil {
		if hasWidthChar(c.message) {
			c.enc = UCS2
		} else {
			c.enc = ASCII
		}
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

// 判断字符串是否包含中文
func hasWidthChar(content string) bool {
	if content == "" {
		return false
	}
	for _, c := range content {
		if c > 0x7f {
			return true
		}
	}
	return false
}

// private static boolean haswidthChar(String content) {
// 	if (StringUtils.isEmpty(content))
// 		return false;

// 	byte[] bytes = content.getBytes();
// 	for (int i = 0; i < bytes.length; i++) {
// 		// 判断最高位是否为1
// 		if ((bytes[i] & (byte) 0x80) == (byte) 0x80) {
// 			return true;
// 		}
// 	}
// 	return false;
// }
// }
