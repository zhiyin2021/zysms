package smpp

import (
	"sync/atomic"

	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/smserror"
)

var (
	ref = uint32(0)
)

// ShortMessage message.
type ShortMessage struct {
	SmDefaultMsgID    byte
	message           string
	enc               codec.Encoding
	udHeader          UDH
	messageData       []byte
	withoutDataCoding bool // purpose of ReplaceSM usage
}

func (c *ShortMessage) MsgLength() int {
	n := len(c.messageData)
	if c.udHeader != nil {
		n += 6
	}
	return n
}

// NewShortMessage returns new ShortMessage.
func NewShortMessage(message string) (s ShortMessage, err error) {
	err = s.SetMessageWithEncoding(message, codec.GSM7BIT)
	return
}

// NewShortMessageWithEncoding returns new ShortMessage with predefined encoding.
func NewShortMessageWithEncoding(message string, enc codec.Encoding) (s ShortMessage, err error) {
	err = s.SetMessageWithEncoding(message, enc)
	return
}

// NewBinaryShortMessage returns new ShortMessage.
func NewBinaryShortMessage(messageData []byte) (s ShortMessage, err error) {
	err = s.SetMessageDataWithEncoding(messageData, codec.BINARY8BIT2)
	return
}

// NewBinaryShortMessageWithEncoding returns new ShortMessage with predefined encoding.
func NewBinaryShortMessageWithEncoding(messageData []byte, enc codec.Encoding) (s ShortMessage, err error) {
	err = s.SetMessageDataWithEncoding(messageData, enc)
	return
}

// NewLongMessage returns long message splitted into multiple short message
func NewLongMessage(message string) (s []*ShortMessage, err error) {
	return NewLongMessageWithEncoding(message, codec.GSM7BIT)
}

// NewLongMessageWithEncoding returns long message splitted into multiple short message with encoding of choice
func NewLongMessageWithEncoding(message string, enc codec.Encoding) (s []*ShortMessage, err error) {
	sm := &ShortMessage{
		message: message,
		enc:     enc,
	}
	return sm.split()
}

// SetMessageWithEncoding sets message with encoding.
func (c *ShortMessage) SetMessageWithEncoding(message string, enc codec.Encoding) (err error) {
	if enc == nil {
		if hasWidthChar(message) {
			if c.messageData, err = codec.UCS2.Encode(message); err == nil {
				c.message = message
				c.enc = codec.UCS2
				return
			}
		}
		if c.messageData, err = codec.GSM7BIT.Encode(message); err == nil {
			c.message = message
			c.enc = codec.GSM7BIT
		} else if c.messageData, err = codec.ASCII.Encode(message); err == nil {
			c.message = message
			c.enc = codec.ASCII
		}
	} else if c.messageData, err = enc.Encode(message); err == nil {
		if len(c.messageData) > SM_MSG_LEN {
			err = smserror.ErrShortMessageLengthTooLarge
		} else {
			c.message = message
			c.enc = enc
		}
	}
	return
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

// SetLongMessageWithEnc sets ShortMessage with message longer than  256 bytes
// callers are expected to call Split() after this
func (c *ShortMessage) SetLongMessageWithEnc(message string, enc codec.Encoding) (err error) {
	c.message = message
	c.enc = enc
	return
}

// UDH gets user data header for short message
func (c *ShortMessage) UDH() UDH {
	return c.udHeader
}

// SetUDH sets user data header for short message
// also appends udh to the beginning of messageData
func (c *ShortMessage) SetUDH(udh UDH) {
	c.udHeader = udh
}

// SetMessageDataWithEncoding sets underlying raw data which is used for pdu marshalling.
func (c *ShortMessage) SetMessageDataWithEncoding(d []byte, enc codec.Encoding) (err error) {
	if len(d) > SM_MSG_LEN {
		err = smserror.ErrShortMessageLengthTooLarge
	} else {
		c.messageData = d
		c.enc = enc
	}
	return
}

// GetMessageData returns underlying binary message.
func (c *ShortMessage) GetMessageData() (d []byte, err error) {
	return c.messageData, nil
}

// GetMessage returns underlying message.
func (c *ShortMessage) GetMessage() (st string, err error) {
	enc := c.enc
	if enc == nil {
		enc = codec.GSM7BIT
	}
	st, err = c.GetMessageWithEncoding(enc)
	return
}

// GetMessageWithEncoding returns (decoded) underlying message.
func (c *ShortMessage) GetMessageWithEncoding(enc codec.Encoding) (st string, err error) {
	if len(c.messageData) > 0 {
		st, err = enc.Decode(c.messageData)
	}
	return
}

// split one short message and split into multiple short message, with UDH
// according to 33GP TS 23.040 section 9.2.3.24.1
//
// NOTE: split() will return array of length 1 if data length is still within the limit
// The encoding interface can implement the Splitter interface for ad-hoc splitting rule
func (c *ShortMessage) split() (multiSM []*ShortMessage, err error) {
	var encoding codec.Encoding

	if c.enc == nil {
		encoding = codec.GSM7BIT
	} else {
		encoding = c.enc
	}

	// check if encoding implements Splitter
	splitter, ok := encoding.(codec.Splitter)
	// check if encoding implements Splitter or split is necessary
	if !ok || !splitter.ShouldSplit(c.message, SM_GSM_MSG_LEN) {
		err = c.SetMessageWithEncoding(c.message, c.enc)
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
	ref := getRefNum()

	// construct SM(s)
	for i, seg := range segments {
		// create new SM, encode data
		multiSM = append(multiSM, &ShortMessage{
			enc: c.enc,
			// message: we don't really care
			messageData:       seg,
			withoutDataCoding: c.withoutDataCoding,
			udHeader:          UDH{NewIEConcatMessage(uint8(len(segments)), uint8(i+1), uint8(ref))},
		})
	}

	return
}

// Marshal implements PDU interface.
func (c *ShortMessage) Marshal(b *codec.BytesWriter) {
	var (
		udhBin []byte
		n      = byte(len(c.messageData))
	)

	// Prepend UDH to message data if there are any
	if c.udHeader != nil && c.udHeader.UDHL() > 0 {
		udhBin, _ = c.udHeader.MarshalBinary()
	}

	b.Grow(int(n) + 3)

	var coding byte
	if c.enc == nil {
		coding = codec.GSM7BITCoding
	} else {
		coding = c.enc.DataCoding()
	}

	// data_coding
	if !c.withoutDataCoding {
		_ = b.WriteByte(coding)
	}

	// sm_default_msg_id
	_ = b.WriteByte(c.SmDefaultMsgID)

	// sm_length
	if udhBin != nil {
		_ = b.WriteByte(byte(int(n) + len(udhBin)))
		b.Write(udhBin)
	} else {
		_ = b.WriteByte(n)
	}

	// short_message
	_, _ = b.Write(c.messageData[:n])
}

// Unmarshal implements PDU interface.
func (c *ShortMessage) Unmarshal(b *codec.BytesReader, udhi bool) (err error) {
	var dataCoding byte
	if !c.withoutDataCoding {
		dataCoding = b.ReadByte()
	}
	c.SmDefaultMsgID = b.ReadByte()
	n := b.ReadByte()
	c.messageData = b.ReadN(int(n))
	if b.Err() != nil {
		return
	}

	c.enc = codec.GetCodec(dataCoding)

	// If short message length is non zero, short message contains User-Data Header
	// Else UDH should be in TLV field MessagePayload
	if udhi && n > 0 {
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

	return
}

// Encoding returns message encoding.
func (c *ShortMessage) Encoding() codec.Encoding {
	return c.enc
}

// returns an atomically incrementing number each time it's called
func getRefNum() uint32 {
	return atomic.AddUint32(&ref, 1)
}

// NOTE:
// When coding splitting function, I have 4 choices of abstraction
// 1. Split the message before encode
// 2. Split the message after encoded
// 3. Split the message DURING encoding (before bit packing)
// 4. Encode, unpack, split
//
// Disadvantages:
// 1. The only way to really know if each segment will fit into 134 octet limit is
//		to do some kind of simulated encoding, where you calculate the total octet
//		by iterating through each character one by one.
//		Too cumbersome
//
// 2. When breaking string at octet position 134, I have to detemeine which
//		character is it ( by doing some kind of decoding)

//		a. If the character code point does not fit in the octet
//		boundary, it has to be carried-over to the next segment.
//		The remaining bits after extracting the carry-over
//		has to be filled with zero.

//		b. If it is an escape character, then I have to backtrack
//		even further since escape chars are not allowed to be splitted
//		in the middle.
//		Since the second bytes of escape chars can be confused with
//		normal chars, I must always lookback 2 character ( repeat step a for at least 2 septet )

//		c. After extracting the carry-on
//		-> Option 2 is very hard when bit packing is already applied
//
// 3. Options 3 require extending Encoding interface,
//	The not good point is not being able to utilize the encoder's Transform() method
//	The good point is you don't have to do bit packing twice

// 4. Terrible option

// All this headaches really only apply to variable length encoding.
// When using fixed length encoding, you can really split the source message BEFORE encodes.
