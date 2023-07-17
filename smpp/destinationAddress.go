package smpp

import (
	"fmt"

	"github.com/zhiyin2021/zysms/codec"
)

// DestinationAddress represents Address or Distribution List based on destination flag.
type DestinationAddress struct {
	destFlag byte
	address  Address
	dl       DistributionList
}

// NewDestinationAddress returns new DestinationAddress.
func NewDestinationAddress() (c DestinationAddress) {
	c.destFlag = DFLT_DEST_FLAG
	return
}

// Unmarshal from buffer.
func (c *DestinationAddress) Unmarshal(b *codec.BytesReader) (err error) {
	if c.destFlag = b.ReadByte(); b.Err() == nil {
		switch c.destFlag {

		case SM_DEST_SME_ADDRESS:
			err = c.address.Unmarshal(b)

		case SM_DEST_DL_NAME:
			err = c.dl.Unmarshal(b)

		default:
			err = fmt.Errorf("unrecognize dest_flag %d", c.destFlag)

		}
	}
	return
}

// Marshal to buffer.
func (c *DestinationAddress) Marshal(b *codec.BytesWriter) {
	switch c.destFlag {
	case SM_DEST_DL_NAME:
		_ = b.WriteByte(SM_DEST_DL_NAME)
		c.dl.Marshal(b)

	default:
		_ = b.WriteByte(SM_DEST_SME_ADDRESS)
		c.address.Marshal(b)
	}
}

// Address returns underlying Address.
func (c *DestinationAddress) Address() Address {
	return c.address
}

// DistributionList returns underlying DistributionList.
func (c *DestinationAddress) DistributionList() DistributionList {
	return c.dl
}

// SetAddress marks DistributionAddress as a SME Address and assign.
func (c *DestinationAddress) SetAddress(addr Address) {
	c.destFlag = SM_DEST_SME_ADDRESS
	c.address = addr
}

// SetDistributionList marks DistributionAddress as a DistributionList and assign.
func (c *DestinationAddress) SetDistributionList(list DistributionList) {
	c.destFlag = SM_DEST_DL_NAME
	c.dl = list
}

// HasValue returns true if underlying DistributionList/Address is assigned.
func (c *DestinationAddress) HasValue() bool {
	return c.destFlag != DFLT_DEST_FLAG
}

// IsAddress returns true if DestinationAddress is a SME Address.
func (c *DestinationAddress) IsAddress() bool {
	return c.destFlag == SM_DEST_SME_ADDRESS
}

// IsDistributionList returns true if DestinationAddress is a DistributionList.
func (c *DestinationAddress) IsDistributionList() bool {
	return c.destFlag == byte(SM_DEST_DL_NAME)
}

// DestinationAddresses represents list of DestinationAddress.
type DestinationAddresses struct {
	l []DestinationAddress
}

// NewDestinationAddresses returns list of DestinationAddress.
func NewDestinationAddresses() (u DestinationAddresses) {
	u.l = make([]DestinationAddress, 0, 8)
	return
}

// Add to list.
func (c *DestinationAddresses) Add(addresses ...DestinationAddress) {
	c.l = append(c.l, addresses...)
}

// Get list.
func (c *DestinationAddresses) Get() []DestinationAddress {
	return c.l
}

// Unmarshal from buffer.
func (c *DestinationAddresses) Unmarshal(b *codec.BytesReader) (err error) {
	var n byte
	if n = b.ReadByte(); b.Err() == nil {
		c.l = make([]DestinationAddress, n)
		var i byte
		for ; i < n; i++ {
			if err = c.l[i].Unmarshal(b); err != nil {
				return
			}
		}
	}
	return
}

// Marshal to buffer.
func (c *DestinationAddresses) Marshal(b *codec.BytesWriter) {
	n := byte(len(c.l))
	b.WriteByte(n)

	var i byte
	for ; i < n; i++ {
		c.l[i].Marshal(b)
	}
}
