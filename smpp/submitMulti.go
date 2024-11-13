package smpp

import "github.com/zhiyin2021/zysms/codec"

// SubmitMulti PDU is used to submit an SMPP message for delivery to multiple recipients
// or to one or more Distribution Lists. The submit_multi PDU does not support
// the transaction message mode.
type SubmitMulti struct {
	base
	ServiceType          string
	SourceAddr           Address
	DestAddrs            DestinationAddresses
	EsmClass             byte
	ProtocolID           byte
	PriorityFlag         byte
	ScheduleDeliveryTime string
	ValidityPeriod       string // not used
	RegisteredDelivery   byte
	ReplaceIfPresentFlag byte // not used
	Message              ShortMessage
}

// NewSubmitMulti returns NewSubmitMulti PDU.
func NewSubmitMulti() codec.PDU {
	message, _ := NewShortMessage("")
	c := &SubmitMulti{
		base:                 newBase(SUBMIT_MULTI, 0),
		ServiceType:          DFLT_SRVTYPE,
		SourceAddr:           NewAddress(),
		DestAddrs:            NewDestinationAddresses(),
		EsmClass:             DFLT_ESM_CLASS,
		ProtocolID:           DFLT_PROTOCOLID,
		PriorityFlag:         DFLT_PRIORITY_FLAG,
		ScheduleDeliveryTime: DFLT_SCHEDULE,
		ValidityPeriod:       DFLT_VALIDITY,
		RegisteredDelivery:   DFLT_REG_DELIVERY,
		ReplaceIfPresentFlag: DFTL_REPLACE_IFP,
		Message:              message,
	}
	return c
}

// GetResponse implements PDU interface.
func (c *SubmitMulti) GetResponse() codec.PDU {
	return &SubmitMultiResp{
		base:          newBase(SUBMIT_MULTI_RESP, c.SequenceNumber),
		MessageID:     DFLT_MSGID,
		UnsuccessSMEs: NewUnsuccessSMEs(),
	}
}

// Marshal implements PDU interface.
func (c *SubmitMulti) Marshal(b *codec.BytesWriter) {
	c.base.marshal(b, func(b *codec.BytesWriter) {
		b.Grow(len(c.ServiceType) + len(c.ScheduleDeliveryTime) + len(c.ValidityPeriod) + 10)

		_ = b.WriteCStr(c.ServiceType)
		c.SourceAddr.Marshal(b)
		c.DestAddrs.Marshal(b)
		_ = b.WriteByte(c.EsmClass)
		_ = b.WriteByte(c.ProtocolID)
		_ = b.WriteByte(c.PriorityFlag)
		_ = b.WriteCStr(c.ScheduleDeliveryTime)
		_ = b.WriteCStr(c.ValidityPeriod)
		_ = b.WriteByte(c.RegisteredDelivery)
		_ = b.WriteByte(c.ReplaceIfPresentFlag)
		c.Message.Marshal(b)
	})
}

// Unmarshal implements PDU interface.
func (c *SubmitMulti) Unmarshal(b *codec.BytesReader) error {
	return c.base.unmarshal(b, func(b *codec.BytesReader) error {
		c.ServiceType = b.ReadCStr()
		c.SourceAddr.Unmarshal(b)
		c.DestAddrs.Unmarshal(b)
		c.EsmClass = b.ReadU8()
		c.ProtocolID = b.ReadU8()
		c.PriorityFlag = b.ReadU8()
		c.ScheduleDeliveryTime = b.ReadCStr()
		c.ValidityPeriod = b.ReadCStr()
		c.RegisteredDelivery = b.ReadU8()
		c.ReplaceIfPresentFlag = b.ReadU8()
		c.Message.Unmarshal(b, (c.EsmClass&SM_UDH_GSM) > 0)
		return b.Err()
	})
}

// SubmitMultiResp PDU.
type SubmitMultiResp struct {
	base
	MessageID     string
	UnsuccessSMEs UnsuccessSMEs
}

// NewSubmitMultiResp returns new SubmitMultiResp.
func NewSubmitMultiResp() codec.PDU {
	c := &SubmitMultiResp{
		base:          newBase(SUBMIT_MULTI_RESP, 0),
		MessageID:     DFLT_MSGID,
		UnsuccessSMEs: NewUnsuccessSMEs(),
	}
	return c
}

// GetResponse implements PDU interface.
func (c *SubmitMultiResp) GetResponse() codec.PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *SubmitMultiResp) Marshal(b *codec.BytesWriter) {
	c.base.marshal(b, func(b *codec.BytesWriter) {
		b.Grow(len(c.MessageID) + 1)
		_ = b.WriteCStr(c.MessageID)
		c.UnsuccessSMEs.Marshal(b)
	})
}

// Unmarshal implements PDU interface.
func (c *SubmitMultiResp) Unmarshal(b *codec.BytesReader) error {
	return c.base.unmarshal(b, func(b *codec.BytesReader) (err error) {
		c.MessageID = b.ReadCStr()
		c.UnsuccessSMEs.Unmarshal(b)
		return b.Err()
	})
}
