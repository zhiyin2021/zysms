package smpp

import (
	"github.com/zhiyin2021/zysms/codec"
)

// SubmitSM PDU is used by an ESME to submit a short message to the SMSC for onward
// transmission to a specified short message entity (SME). The submit_sm PDU does
// not support the transaction message mode.
type SubmitSM struct {
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
	MessageID            string
}

// NewSubmitSM returns SubmitSM PDU.
func NewSubmitSM() codec.PDU {
	message, _ := NewShortMessage("")
	c := &SubmitSM{
		base:                 newBase(SUBMIT_SM, 0),
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

// ShouldSplit check if this the user data of submitSM PDU
func (c *SubmitSM) ShouldSplit() bool {
	// GSM standard mandates that User Data must be no longer than 140 octet
	return len(c.Message.messageData) > 140
}

// GetResponse implements PDU interface.
func (c *SubmitSM) GetResponse() codec.PDU {
	return &SubmitSMResp{
		base:      newBase(SUBMIT_SM_RESP, c.SequenceNumber),
		MessageID: DFLT_MSGID,
	}
}

// Split split a single long text message into multiple SubmitSM PDU,
// Each have the TPUD within the GSM's User Data limit of 140 octet
// If the message is short enough and doesn't need splitting,
// Split() returns an array of length 1
func (c *SubmitSM) Split() (multiSubSM []*SubmitSM, err error) {
	multiSubSM = []*SubmitSM{}

	multiMsg, err := c.Message.split()
	if err != nil {
		return
	}

	for _, msg := range multiMsg {
		multiSubSM = append(multiSubSM, &SubmitSM{
			base:                 c.base,
			ServiceType:          c.ServiceType,
			SourceAddr:           c.SourceAddr,
			DestAddr:             c.DestAddr,
			EsmClass:             c.EsmClass | SM_UDH_GSM, // must set to indicate UDH
			ProtocolID:           c.ProtocolID,
			PriorityFlag:         c.PriorityFlag,
			ScheduleDeliveryTime: c.ScheduleDeliveryTime,
			ValidityPeriod:       c.ValidityPeriod,
			RegisteredDelivery:   c.RegisteredDelivery,
			ReplaceIfPresentFlag: c.ReplaceIfPresentFlag,
			Message:              *msg,
		})
	}
	return
}

// Marshal implements PDU interface.
func (c *SubmitSM) Marshal(b *codec.BytesWriter) {
	c.base.marshal(b, func(b *codec.BytesWriter) {
		b.Grow(len(c.ServiceType) + len(c.ScheduleDeliveryTime) + len(c.ValidityPeriod) + 10)

		b.WriteCStr(c.ServiceType)
		c.SourceAddr.Marshal(b)
		c.DestAddr.Marshal(b)
		b.WriteByte(c.EsmClass)
		b.WriteByte(c.ProtocolID)
		b.WriteByte(c.PriorityFlag)
		b.WriteCStr(c.ScheduleDeliveryTime)
		b.WriteCStr(c.ValidityPeriod)
		b.WriteByte(c.RegisteredDelivery)
		b.WriteByte(c.ReplaceIfPresentFlag)
		c.Message.Marshal(b)
	})
}

// Unmarshal implements PDU interface.
func (c *SubmitSM) Unmarshal(b *codec.BytesReader) error {
	return c.base.unmarshal(b, func(b *codec.BytesReader) (err error) {
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

// SubmitSMResp PDU.
type SubmitSMResp struct {
	base
	MessageID string
}

// NewSubmitSMResp returns new SubmitSMResp.
func NewSubmitSMResp() codec.PDU {
	c := &SubmitSMResp{
		base:      newBase(SUBMIT_SM_RESP, 0),
		MessageID: DFLT_MSGID,
	}
	return c
}

// GetResponse implements PDU interface.
func (c *SubmitSMResp) GetResponse() codec.PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *SubmitSMResp) Marshal(b *codec.BytesWriter) {
	c.base.marshal(b, func(b *codec.BytesWriter) {
		if c.CommandStatus == ESME_ROK {
			b.Grow(len(c.MessageID) + 1)
			_ = b.WriteCStr(c.MessageID)
		}
	})
}

// Unmarshal implements PDU interface.
func (c *SubmitSMResp) Unmarshal(b *codec.BytesReader) error {
	return c.base.unmarshal(b, func(b *codec.BytesReader) (err error) {
		// if c.CommandStatus == ESME_ROK {
		c.MessageID = b.ReadCStr()
		// }
		return nil
	})
}
