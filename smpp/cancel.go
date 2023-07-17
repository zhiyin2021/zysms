package smpp

import "github.com/zhiyin2021/zysms/codec"

// CancelSM PDU is issued by the ESME to cancel one or more previously submitted short messages
// that are still pending delivery. The command may specify a particular message to cancel, or
// all messages for a particular source, destination and service_type are to be cancelled.
type CancelSM struct {
	base
	ServiceType string
	MessageID   string
	SourceAddr  Address
	DestAddr    Address
}

// NewCancelSM returns CancelSM PDU.
func NewCancelSM() codec.PDU {
	c := &CancelSM{
		base:        newBase(CANCEL_SM, 0),
		ServiceType: DFLT_SRVTYPE,
		MessageID:   DFLT_MSGID,
		SourceAddr:  NewAddress(),
		DestAddr:    NewAddress(),
	}
	return c
}

// GetResponse implements PDU interface.
func (c *CancelSM) GetResponse() codec.PDU {
	return &CancelSMResp{
		base: newBase(CANCEL_SM_RESP, c.SequenceNumber),
	}
}

// Marshal implements PDU interface.
func (c *CancelSM) Marshal(b *codec.BytesWriter) {
	c.base.marshal(b, func(b *codec.BytesWriter) {
		b.Grow(len(c.ServiceType) + len(c.MessageID) + 2)
		b.WriteCStr(c.ServiceType)
		b.WriteCStr(c.MessageID)
		c.SourceAddr.Marshal(b)
		c.DestAddr.Marshal(b)
	})
}

// Unmarshal implements PDU interface.
func (c *CancelSM) Unmarshal(b *codec.BytesReader) error {
	return c.base.unmarshal(b, func(b *codec.BytesReader) error {
		c.ServiceType = b.ReadCStr()
		c.MessageID = b.ReadCStr()
		c.SourceAddr.Unmarshal(b)
		c.DestAddr.Unmarshal(b)
		return b.Err()
	})
}

// CancelSMResp PDU.
type CancelSMResp struct {
	base
}

// NewCancelSMResp returns CancelSMResp.
func NewCancelSMResp() codec.PDU {
	c := &CancelSMResp{
		base: newBase(CANCEL_SM_RESP, 0),
	}
	return c
}

// GetResponse implements PDU interface.
func (c *CancelSMResp) GetResponse() codec.PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *CancelSMResp) Marshal(b *codec.BytesWriter) {
	c.base.marshal(b, nil)
}

// Unmarshal implements PDU interface.
func (c *CancelSMResp) Unmarshal(b *codec.BytesReader) error {
	return c.base.unmarshal(b, nil)
}
