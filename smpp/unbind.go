package smpp

import (
	"github.com/zhiyin2021/zysms/codec"
)

// Unbind PDU is to deregister an instance of an ESME from the SMSC and inform the SMSC
// that the ESME no longer wishes to use this network connection for the submission or
// delivery of messages.
type Unbind struct {
	base
}

// NewUnbind returns Unbind PDU.
func NewUnbind() codec.PDU {
	c := &Unbind{
		base: newBase(UNBIND, 0),
	}
	return c
}

// GetResponse implements PDU interface.
func (c *Unbind) GetResponse() codec.PDU {
	return &UnbindResp{
		base: newBase(UNBIND_RESP, c.SequenceNumber),
	}
}

// Marshal implements PDU interface.
func (c *Unbind) Marshal(b *codec.BytesWriter) {
	c.base.marshal(b, nil)
}

// Unmarshal implements PDU interface.
func (c *Unbind) Unmarshal(b *codec.BytesReader) error {
	return c.base.unmarshal(b, nil)
}

// UnbindResp PDU.
type UnbindResp struct {
	base
}

// NewUnbindResp returns UnbindResp.
func NewUnbindResp() codec.PDU {
	c := &UnbindResp{
		base: newBase(UNBIND_RESP, 0),
	}
	return c
}

// GetResponse implements PDU interface.
func (c *UnbindResp) GetResponse() codec.PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *UnbindResp) Marshal(b *codec.BytesWriter) {
	c.base.marshal(b, nil)
}

// Unmarshal implements PDU interface.
func (c *UnbindResp) Unmarshal(b *codec.BytesReader) error {
	return c.base.unmarshal(b, nil)
}
