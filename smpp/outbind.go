package smpp

import "github.com/zhiyin2021/zysms/codec"

// Outbind PDU is used by the SMSC to signal an ESME to originate a bind_receiver request to the SMSC.
type Outbind struct {
	base
	SystemID string
	Password string
}

// NewOutbind returns Outbind PDU.
func NewOutbind() codec.PDU {
	c := &Outbind{
		base: newBase(OUTBIND, 0),
	}
	return c
}

// GetResponse implements PDU interface.
func (c *Outbind) GetResponse() codec.PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *Outbind) Marshal(b *codec.BytesWriter) {
	c.base.marshal(b, func(b *codec.BytesWriter) {
		b.Grow(len(c.SystemID) + len(c.Password) + 2)
		b.WriteCStr(c.SystemID)
		b.WriteCStr(c.Password)
	})
}

// Unmarshal implements PDU interface.
func (c *Outbind) Unmarshal(b *codec.BytesReader) error {
	return c.base.unmarshal(b, func(b *codec.BytesReader) (err error) {
		c.SystemID = b.ReadCStr()
		c.Password = b.ReadCStr()
		return b.Err()
	})
}
