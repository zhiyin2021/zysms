package smpp

import (
	"github.com/zhiyin2021/zysms/codec"
)

// UnsuccessSME indicates submission was unsuccessful and the respective errors.
type UnsuccessSME struct {
	Address
	errorStatusCode CommandStatus
}

// NewUnsuccessSME returns new UnsuccessSME
func NewUnsuccessSME() (c UnsuccessSME) {
	c = UnsuccessSME{
		Address:         NewAddress(),
		errorStatusCode: ESME_ROK,
	}
	return
}

// NewUnsuccessSMEWithAddr returns new UnsuccessSME with address.
func NewUnsuccessSMEWithAddr(addr string, status codec.CommandStatus) (c UnsuccessSME, err error) {
	c = NewUnsuccessSME()
	if err = c.SetAddress(addr); err == nil {
		c.SetErrorStatusCode(status)
	}
	return
}

// NewUnsuccessSMEWithTonNpi create new address with ton, npi and error code.
func NewUnsuccessSMEWithTonNpi(ton, npi byte, status codec.CommandStatus) UnsuccessSME {
	return UnsuccessSME{
		Address:         NewAddressWithTonNpi(ton, npi),
		errorStatusCode: status,
	}
}

// Unmarshal from buffer.
func (c *UnsuccessSME) Unmarshal(b *codec.BytesReader) (err error) {
	if err = c.Address.Unmarshal(b); err == nil {
		c.errorStatusCode = codec.CommandStatus(b.ReadU32())
	}
	return
}

// Marshal to buffer.
func (c *UnsuccessSME) Marshal(b *codec.BytesWriter) {
	c.Address.Marshal(b)
	b.WriteU32(uint32(c.errorStatusCode))
}

// SetErrorStatusCode sets error status code.
func (c *UnsuccessSME) SetErrorStatusCode(v codec.CommandStatus) {
	c.errorStatusCode = v
}

// ErrorStatusCode returns assigned status code.
func (c *UnsuccessSME) ErrorStatusCode() codec.CommandStatus {
	return c.errorStatusCode
}

// UnsuccessSMEs represents list of UnsuccessSME.
type UnsuccessSMEs struct {
	l []UnsuccessSME
}

// NewUnsuccessSMEs returns list of UnsuccessSME.
func NewUnsuccessSMEs() (u UnsuccessSMEs) {
	u.l = make([]UnsuccessSME, 0, 8)
	return
}

// Add to list.
func (c *UnsuccessSMEs) Add(us ...UnsuccessSME) {
	c.l = append(c.l, us...)
}

// Get list.
func (c *UnsuccessSMEs) Get() []UnsuccessSME {
	return c.l
}

// Unmarshal from buffer.
func (c *UnsuccessSMEs) Unmarshal(b *codec.BytesReader) (err error) {
	n := b.ReadByte()
	if err = b.Err(); err == nil {
		c.l = make([]UnsuccessSME, n)
		var i byte
		for ; i < n; i++ {
			if err = c.l[i].Unmarshal(b); err != nil {
				return
			}
		}
	}
	return err
}

// Marshal to buffer.
func (c *UnsuccessSMEs) Marshal(b *codec.BytesWriter) {
	n := byte(len(c.l))
	_ = b.WriteByte(n)

	var i byte
	for ; i < n; i++ {
		c.l[i].Marshal(b)
	}
}
