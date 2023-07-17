package smpp

import "github.com/zhiyin2021/zysms/codec"

// EnquireLink PDU. This message can be sent by either the ESME or SMSC
// and is used to provide a confidence- check of the communication path between
// an ESME and an SMSC. On receipt of this request the receiving party should
// respond with an enquire_link_resp, thus verifying that the application
// level connection between the SMSC and the ESME is functioning.
// The ESME may also respond by sending any valid SMPP primitive.
type EnquireLink struct {
	base
}

// NewEnquireLink returns new EnquireLink PDU.
func NewEnquireLink() codec.PDU {
	c := &EnquireLink{
		base: newBase(ENQUIRE_LINK, 0),
	}
	return c
}

// GetResponse implements PDU interface.
func (c *EnquireLink) GetResponse() codec.PDU {
	return &EnquireLinkResp{
		base: newBase(ENQUIRE_LINK_RESP, c.SequenceNumber),
	}
}

// Marshal implements PDU interface.
func (c *EnquireLink) Marshal(b *codec.BytesWriter) {
	c.base.marshal(b, nil)
}

// Unmarshal implements PDU interface.
func (c *EnquireLink) Unmarshal(b *codec.BytesReader) error {
	return c.base.unmarshal(b, nil)
}

// EnquireLinkResp PDU.
type EnquireLinkResp struct {
	base
}

// NewEnquireLinkResp returns EnquireLinkResp.
func NewEnquireLinkResp() codec.PDU {
	c := &EnquireLinkResp{
		base: newBase(ENQUIRE_LINK_RESP, 0),
	}
	return c
}

// GetResponse implements PDU interface.
func (c *EnquireLinkResp) GetResponse() codec.PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *EnquireLinkResp) Marshal(b *codec.BytesWriter) {
	c.base.marshal(b, nil)
}

// Unmarshal implements PDU interface.
func (c *EnquireLinkResp) Unmarshal(b *codec.BytesReader) error {
	return c.base.unmarshal(b, nil)
}
