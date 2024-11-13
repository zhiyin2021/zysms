package smpp

import "github.com/zhiyin2021/zysms/codec"

// ReplaceSM PDU is issued by the ESME to replace a previously submitted short message
// that is still pending delivery. The matching mechanism is based on the message_id and
// source address of the original message. Where the original submit_sm ‘source address’
// was defaulted to NULL, then the source address in the replace_sm command should also be NULL.
type ReplaceSM struct {
	base
	MessageID            string
	SourceAddr           Address
	ScheduleDeliveryTime string
	ValidityPeriod       string
	RegisteredDelivery   byte
	Message              ShortMessage
}

// NewReplaceSM returns ReplaceSM PDU.
func NewReplaceSM() codec.PDU {
	message, _ := NewShortMessage("")
	message.withoutDataCoding = true
	c := &ReplaceSM{
		base:                 newBase(REPLACE_SM, 0),
		SourceAddr:           NewAddress(),
		ScheduleDeliveryTime: DFLT_SCHEDULE,
		ValidityPeriod:       DFLT_VALIDITY,
		RegisteredDelivery:   DFLT_REG_DELIVERY,
		Message:              message,
	}
	return c
}

// GetResponse implements PDU interface.
func (c *ReplaceSM) GetResponse() codec.PDU {
	return &ReplaceSMResp{
		base: newBase(REPLACE_SM_RESP, c.SequenceNumber),
	}
}

// Marshal implements PDU interface.
func (c *ReplaceSM) Marshal(b *codec.BytesWriter) {
	c.base.marshal(b, func(b *codec.BytesWriter) {
		b.Grow(len(c.MessageID) + len(c.ScheduleDeliveryTime) + len(c.ValidityPeriod) + 4)

		_ = b.WriteCStr(c.MessageID)
		c.SourceAddr.Marshal(b)
		_ = b.WriteCStr(c.ScheduleDeliveryTime)
		_ = b.WriteCStr(c.ValidityPeriod)
		_ = b.WriteByte(c.RegisteredDelivery)
		c.Message.Marshal(b)
	})
}

// Unmarshal implements PDU interface.
func (c *ReplaceSM) Unmarshal(b *codec.BytesReader) error {
	return c.base.unmarshal(b, func(b *codec.BytesReader) error {
		c.MessageID = b.ReadCStr()
		c.SourceAddr.Unmarshal(b)
		c.ScheduleDeliveryTime = b.ReadCStr()
		c.ValidityPeriod = b.ReadCStr()
		c.RegisteredDelivery = b.ReadU8()
		c.Message.Unmarshal(b, false)
		return b.Err()
	})
}

// ReplaceSMResp PDU.
type ReplaceSMResp struct {
	base
}

// NewReplaceSMResp returns ReplaceSMResp.
func NewReplaceSMResp() codec.PDU {
	c := &ReplaceSMResp{
		base: newBase(REPLACE_SM_RESP, 0),
	}
	return c
}

// GetResponse implements PDU interface.
func (c *ReplaceSMResp) GetResponse() codec.PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *ReplaceSMResp) Marshal(b *codec.BytesWriter) {
	c.base.marshal(b, nil)
}

// Unmarshal implements PDU interface.
func (c *ReplaceSMResp) Unmarshal(b *codec.BytesReader) error {
	return c.base.unmarshal(b, nil)
}
