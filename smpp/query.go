package smpp

import "github.com/zhiyin2021/zysms/codec"

// QuerySM PDU is issued by the ESME to query the status of a previously submitted short message.
// The matching mechanism is based on the SMSC assigned message_id and source address. Where the
// original submit_sm, data_sm or submit_multi ‘source address’ was defaulted to NULL, then the
// source address in the query_sm command should also be set to NULL.
type QuerySM struct {
	base
	MessageID  string
	SourceAddr Address
}

// NewQuerySM returns new QuerySM PDU.
func NewQuerySM() codec.PDU {
	c := &QuerySM{
		SourceAddr: NewAddress(),
	}
	c.CommandID = QUERY_SM
	return c
}

// GetResponse implements PDU interface.
func (c *QuerySM) GetResponse() codec.PDU {
	return &QuerySMResp{
		base:         newBase(QUERY_SM_RESP, c.SequenceNumber),
		FinalDate:    DFLT_DATE,
		MessageState: DFLT_MSG_STATE,
		ErrorCode:    DFLT_ERR,
	}
}

// Marshal implements PDU interface.
func (c *QuerySM) Marshal(b *codec.BytesWriter) {
	c.base.marshal(b, func(b *codec.BytesWriter) {
		b.Grow(len(c.MessageID) + 1)

		b.WriteCStr(c.MessageID)
		c.SourceAddr.Marshal(b)
	})
}

// Unmarshal implements PDU interface.
func (c *QuerySM) Unmarshal(b *codec.BytesReader) error {
	return c.base.unmarshal(b, func(b *codec.BytesReader) (err error) {
		c.MessageID = b.ReadCStr()
		c.SourceAddr.Unmarshal(b)
		return b.Err()
	})
}

// QuerySMResp PDU.
type QuerySMResp struct {
	base
	MessageID    string
	FinalDate    string
	MessageState byte
	ErrorCode    byte
}

// NewQuerySMResp returns new QuerySM PDU.
func NewQuerySMResp() codec.PDU {
	c := &QuerySMResp{
		base:         newBase(QUERY_SM_RESP, 0),
		FinalDate:    DFLT_DATE,
		MessageState: DFLT_MSG_STATE,
		ErrorCode:    DFLT_ERR,
	}
	return c
}

// GetResponse implements PDU interface.
func (c *QuerySMResp) GetResponse() codec.PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *QuerySMResp) Marshal(b *codec.BytesWriter) {
	c.base.marshal(b, func(b *codec.BytesWriter) {
		b.Grow(len(c.MessageID) + len(c.FinalDate) + 4)

		b.WriteCStr(c.MessageID)
		b.WriteCStr(c.FinalDate)
		b.WriteByte(c.MessageState)
		b.WriteByte(c.ErrorCode)
	})
}

// Unmarshal implements PDU interface.
func (c *QuerySMResp) Unmarshal(b *codec.BytesReader) error {
	return c.base.unmarshal(b, func(b *codec.BytesReader) error {
		c.MessageID = b.ReadCStr()
		c.FinalDate = b.ReadCStr()
		c.MessageState = b.ReadU8()
		c.ErrorCode = b.ReadU8()
		return b.Err()
	})
}
