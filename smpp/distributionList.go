package smpp

import (
	"fmt"

	"github.com/zhiyin2021/zysms/codec"
)

// DistributionList represents group of contacts.
type DistributionList struct {
	name string
}

// NewDistributionList returns a new DistributionList.
func NewDistributionList(name string) (c DistributionList, err error) {
	err = c.SetName(name)
	return
}

// Unmarshal from buffer.
func (c *DistributionList) Unmarshal(b *codec.BytesReader) error {
	c.name = b.ReadCStr()
	return b.Err()
}

// Marshal to buffer.
func (c *DistributionList) Marshal(b *codec.BytesWriter) {
	b.Grow(1 + len(c.name))
	b.WriteCStr(c.name)
}

// SetName sets DistributionList name.
func (c *DistributionList) SetName(name string) error {
	if len(name) > SM_DL_NAME_LEN {
		return fmt.Errorf("Distribution List name exceed limit. (%d > %d)", len(name), SM_DL_NAME_LEN)
	} else {
		c.name = name
	}
	return nil
}

// Name returns name of DistributionList
func (c DistributionList) Name() string {
	return c.name
}
