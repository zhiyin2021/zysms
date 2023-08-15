package smpp

import "github.com/zhiyin2021/zysms/codec"

const ReportLen byte = 66

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
	Report               *DeliverReport
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
		if c.RegisteredDelivery == 1 && c.Report != nil {
			b.WriteStr(c.Report.MsgId, 10)
			b.WriteStr(c.Report.Sub, 3)
			b.WriteStr(c.Report.Dlvrd, 3)
			b.WriteStr(c.Report.SubmitDate, 10)
			b.WriteStr(c.Report.DoneDate, 10)
			b.WriteStr(c.Report.Stat, 7)
			b.WriteStr(c.Report.Text, 20)
		} else {
			c.Message.Marshal(b)
		}
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
		if c.RegisteredDelivery == 1 {
			c.Report = &DeliverReport{
				MsgId:      b.ReadStr(10),
				Sub:        b.ReadStr(3),
				Dlvrd:      b.ReadStr(3),
				SubmitDate: b.ReadStr(10),
				DoneDate:   b.ReadStr(10),
				Stat:       b.ReadStr(7),
				Text:       b.ReadStr(20),
			}
		} else {
			c.Message.Unmarshal(b, (c.EsmClass&SM_UDH_GSM) > 0)
		}
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

type DeliverReport struct {
	MsgId      string // 10字节 The message ID allocated to the message by the SMSC when originally submitted.
	Sub        string // 3 字节 Number of short messages originally submitted. This is only relevant when the original message was submitted to a distribution list.The value is padded with leading zeros if necessary.
	Dlvrd      string // 3 字节 Number of short messages delivered. This is only relevant where the original message was submitted to a distribution list.The value is padded with leading zeros if necessary.
	SubmitDate string // 10字节 (YYMMDDhhmm)The time and date at which the short message was submitted. In the case of a message which has been replaced, this is the date that the original message was replaced.The format is as follows:
	DoneDate   string // 10字节 (YYMMDDhhmm)The time and date at which the short message reached it’s final state. The format is the same as for the submit date.
	Stat       string // 7 字节 The final status of the message. For settings for this field see Table B-2.
	Text       string // 20字节 The first 20 characters of the short message.
}

/*
DELIVRD
EXPIRED
DELETED
UNDELIV
ACCEPTD
UNKNOWN
REJECTD
*/
