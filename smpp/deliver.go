package smpp

import "github.com/zhiyin2021/zysms/codec"

// DeliverSM PDU is issued by the SMSC to send a message to an ESME.
// Using this command, the SMSC may route a short message to the ESME for delivery.
type DeliverSM struct {
	base
	ServiceType          string
	SourceAddr           Address
	DestAddr             Address
	EsmClass             byte
	ProtocolID           byte
	PriorityFlag         byte
	ScheduleDeliveryTime string // not used
	ValidityPeriod       string // not used
	RegisteredDelivery   byte
	ReplaceIfPresentFlag byte // not used
	Message              ShortMessage
}

// NewDeliverSM returns DeliverSM PDU.
func NewDeliverSM() codec.PDU {
	message, _ := NewShortMessage("")
	c := &DeliverSM{
		base:                 newBase(DELIVER_SM, 0),
		ServiceType:          DFLT_SRVTYPE,
		SourceAddr:           NewAddress(),
		DestAddr:             NewAddress(),
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
func (c *DeliverSM) GetResponse() codec.PDU {
	return &DeliverSMResp{
		base:      newBase(DELIVER_SM_RESP, c.SequenceNumber),
		MessageID: DFLT_MSGID,
	}
}

// Marshal implements PDU interface.
func (c *DeliverSM) Marshal(b *codec.BytesWriter) {
	c.base.marshal(b, func(b *codec.BytesWriter) {
		b.Grow(len(c.ServiceType) + len(c.ScheduleDeliveryTime) + len(c.ValidityPeriod) + 10)

		_ = b.WriteCStr(c.ServiceType)
		c.SourceAddr.Marshal(b)
		c.DestAddr.Marshal(b)
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
func (c *DeliverSM) Unmarshal(b *codec.BytesReader) error {
	return c.base.unmarshal(b, func(b *codec.BytesReader) error {
		c.ServiceType = b.ReadCStr()
		c.SourceAddr.Unmarshal(b)
		c.DestAddr.Unmarshal(b)
		c.EsmClass = b.ReadByte()
		c.ProtocolID = b.ReadByte()
		c.PriorityFlag = b.ReadByte()
		c.ScheduleDeliveryTime = b.ReadCStr()
		c.ValidityPeriod = b.ReadCStr()
		c.RegisteredDelivery = b.ReadByte()
		c.ReplaceIfPresentFlag = b.ReadByte()
		c.Message.Unmarshal(b, (c.EsmClass&SM_UDH_GSM) > 0)
		return b.Err()
	})
}

// DeliverSMResp PDU.
type DeliverSMResp struct {
	base
	MessageID string
}

// NewDeliverSMResp returns new DeliverSMResp.
func NewDeliverSMResp() codec.PDU {
	c := &DeliverSMResp{
		base:      newBase(DELIVER_SM_RESP, 0),
		MessageID: DFLT_MSGID,
	}
	return c
}

// GetResponse implements PDU interface.
func (c *DeliverSMResp) GetResponse() codec.PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *DeliverSMResp) Marshal(b *codec.BytesWriter) {
	c.base.marshal(b, func(b *codec.BytesWriter) {
		b.Grow(len(c.MessageID) + 1)
		b.WriteCStr(c.MessageID)
	})
}

// Unmarshal implements PDU interface.
func (c *DeliverSMResp) Unmarshal(b *codec.BytesReader) error {
	return c.base.unmarshal(b, func(b *codec.BytesReader) (err error) {
		c.MessageID = b.ReadCStr()
		return b.Err()
	})
}
