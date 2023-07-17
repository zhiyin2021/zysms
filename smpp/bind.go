package smpp

import "github.com/zhiyin2021/zysms/codec"

// BindingType indicates type of binding.
type BindingType byte

const (
	// Receiver indicates Receiver binding.
	Receiver BindingType = iota
	// Transceiver indicates Transceiver binding.
	Transceiver
	// Transmitter indicate Transmitter binding.
	Transmitter
)

// BindRequest represents a bind request.
type BindRequest struct {
	base
	SystemID         string
	Password         string
	SystemType       string
	InterfaceVersion codec.Version
	AddressRange     AddressRange
	BindingType      BindingType
}

// NewBindRequest returns new bind request.
func NewBindRequest(t BindingType) (b *BindRequest) {
	b = &BindRequest{
		base:             newBase(BIND_TRANSCEIVER, 0),
		BindingType:      t,
		SystemID:         DFLT_SYSID,
		Password:         DFLT_PASS,
		SystemType:       DFLT_SYSTYPE,
		AddressRange:     NewAddressRange(),
		InterfaceVersion: V34,
	}
	switch t {
	case Transceiver:
		b.CommandID = BIND_TRANSCEIVER
	case Receiver:
		b.CommandID = BIND_RECEIVER
	case Transmitter:
		b.CommandID = BIND_TRANSMITTER
	}
	return
}

// NewBindTransmitter returns new bind transmitter pdu.
func NewBindTransmitter() codec.PDU {
	return NewBindRequest(Transmitter)
}

// NewBindTransceiver returns new bind transceiver pdu.
func NewBindTransceiver() codec.PDU {
	return NewBindRequest(Transceiver)
}

// NewBindReceiver returns new bind receiver pdu.
func NewBindReceiver() codec.PDU {
	return NewBindRequest(Receiver)
}

// GetResponse implements PDU interface.
func (b *BindRequest) GetResponse() codec.PDU {
	c := &BindResp{
		base: newBase(BIND_TRANSCEIVER_RESP, b.SequenceNumber),
	}
	switch b.BindingType {
	case Transceiver:
		c.CommandID = BIND_TRANSCEIVER_RESP

	case Receiver:
		c.CommandID = BIND_RECEIVER_RESP

	case Transmitter:
		c.CommandID = BIND_TRANSMITTER_RESP
	}

	return c
}

// Marshal implements PDU interface.
func (b *BindRequest) Marshal(w *codec.BytesWriter) {
	b.base.marshal(w, func(w *codec.BytesWriter) {
		w.Grow(len(b.SystemID) + len(b.Password) + len(b.SystemType) + 4)

		w.WriteCStr(b.SystemID)
		w.WriteCStr(b.Password)
		w.WriteCStr(b.SystemType)
		w.WriteByte(byte(b.InterfaceVersion))
		b.AddressRange.Marshal(w)
	})
}

// Unmarshal implements PDU interface.
func (b *BindRequest) Unmarshal(w *codec.BytesReader) error {
	return b.base.unmarshal(w, func(w *codec.BytesReader) error {
		b.SystemID = w.ReadCStr()
		b.Password = w.ReadCStr()
		b.SystemType = w.ReadCStr()

		b.InterfaceVersion = codec.Version(w.ReadByte())
		b.AddressRange.Unmarshal(w)

		return w.Err()
	})
}

// BindResp PDU.
type BindResp struct {
	base
	SystemID string
}

// NewBindTransmitterResp returns new bind transmitter resp.
func NewBindTransmitterResp() codec.PDU {
	c := &BindResp{
		base: newBase(BIND_TRANSMITTER_RESP, 0),
	}
	return c
}

// NewBindTransceiverResp returns new bind transceiver resp.
func NewBindTransceiverResp() codec.PDU {
	c := &BindResp{
		base: newBase(BIND_TRANSCEIVER_RESP, 0),
	}
	return c
}

// NewBindReceiverResp returns new bind receiver resp.
func NewBindReceiverResp() codec.PDU {
	c := &BindResp{
		base: newBase(BIND_RECEIVER_RESP, 0),
	}
	return c
}

// GetResponse implements PDU interface.
func (c *BindResp) GetResponse() codec.PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *BindResp) Marshal(b *codec.BytesWriter) {
	c.base.marshal(b, func(w *codec.BytesWriter) {
		w.Grow(len(c.SystemID) + 1)
		w.WriteCStr(c.SystemID)
	})
}

// Unmarshal implements PDU interface.
func (c *BindResp) Unmarshal(b *codec.BytesReader) error {
	return c.base.unmarshal(b, func(w *codec.BytesReader) (err error) {
		if c.CommandID == BIND_TRANSCEIVER_RESP || c.CommandStatus == ESME_ROK {
			c.SystemID = w.ReadCStr()
		}
		return w.Err()
	})
}
