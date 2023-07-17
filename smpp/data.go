package smpp

import "github.com/zhiyin2021/zysms/codec"

// DataSM PDU is used to transfer data between the SMSC and the ESME.
// It may be used by both the ESME and SMSC.
type DataSM struct {
	base
	ServiceType        string
	SourceAddr         Address
	DestAddr           Address
	EsmClass           byte
	RegisteredDelivery byte
	DataCoding         byte
}

// NewDataSM returns new data sm pdu.
func NewDataSM() codec.PDU {
	c := &DataSM{
		base:               newBase(DATA_SM, 0),
		ServiceType:        DFLT_SRVTYPE,
		SourceAddr:         NewAddress(),
		DestAddr:           NewAddress(),
		EsmClass:           DFLT_ESM_CLASS,
		RegisteredDelivery: DFLT_REG_DELIVERY,
		DataCoding:         DFLT_DATA_CODING,
	}
	return c
}

// GetResponse implements PDU interface.
func (c *DataSM) GetResponse() codec.PDU {
	return &DataSMResp{
		base:      newBase(DATA_SM_RESP, c.SequenceNumber),
		MessageID: DFLT_MSGID,
	}
}

// Marshal implements PDU interface.
func (c *DataSM) Marshal(b *codec.BytesWriter) {
	c.base.marshal(b, func(b *codec.BytesWriter) {
		b.Grow(len(c.ServiceType) + 4)

		_ = b.WriteCStr(c.ServiceType)
		c.SourceAddr.Marshal(b)
		c.DestAddr.Marshal(b)
		_ = b.WriteByte(c.EsmClass)
		_ = b.WriteByte(c.RegisteredDelivery)
		_ = b.WriteByte(c.DataCoding)
	})
}

// Unmarshal implements PDU interface.
func (c *DataSM) Unmarshal(b *codec.BytesReader) error {
	return c.base.unmarshal(b, func(b *codec.BytesReader) (err error) {
		c.ServiceType = b.ReadCStr()
		c.SourceAddr.Unmarshal(b)
		c.DestAddr.Unmarshal(b)
		c.EsmClass = b.ReadByte()
		c.RegisteredDelivery = b.ReadByte()
		c.DataCoding = b.ReadByte()
		return b.Err()
	})
}

// DataSMResp PDU.
type DataSMResp struct {
	base
	MessageID string
}

// NewDataSMResp returns DataSMResp.
func NewDataSMResp() codec.PDU {
	c := &DataSMResp{
		base:      newBase(DATA_SM_RESP, 0),
		MessageID: DFLT_MSGID,
	}
	return c
}

// GetResponse implements PDU interface.
func (c *DataSMResp) GetResponse() codec.PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *DataSMResp) Marshal(b *codec.BytesWriter) {
	c.base.marshal(b, func(b *codec.BytesWriter) {
		b.Grow(len(c.MessageID) + 1)
		b.WriteCStr(c.MessageID)
	})
}

// Unmarshal implements PDU interface.
func (c *DataSMResp) Unmarshal(b *codec.BytesReader) error {
	return c.base.unmarshal(b, func(b *codec.BytesReader) error {
		c.MessageID = b.ReadCStr()
		return b.Err()
	})
}
