package smgp

import "github.com/zhiyin2021/zysms/codec"

type ActiveTestReq struct {
	base
}
type ActiveTestResp struct {
	base
}

func NewActiveTestReq(ver Version) codec.PDU {
	return &ActiveTestReq{
		base: newBase(ver, SMGP_ACTIVE_TEST, 0),
	}
}
func NewActiveTestResp(ver Version) codec.PDU {
	return &ActiveTestResp{
		base: newBase(ver, SMGP_ACTIVE_TEST_RESP, 0),
	}
}

// Pack packs the ActiveTestReq to bytes stream for client side.
func (p *ActiveTestReq) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, nil)
}

// Unpack unpack the binary byte stream to a ActiveTestReq variable.
// After unpack, you will get all value of fields in
// ActiveTestReq struct.
func (p *ActiveTestReq) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, nil)
}

// GetResponse implements PDU interface.
func (b *ActiveTestReq) GetResponse() codec.PDU {
	return &ActiveTestResp{
		base: newBase(b.Version, SMGP_ACTIVE_TEST_RESP, b.SequenceNumber),
	}
}

// Pack packs the ActiveTestReq to bytes stream for client side.
func (p *ActiveTestResp) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, nil)
}

// Unpack unpack the binary byte stream to a ActiveTestReq variable.
// After unpack, you will get all value of fields in
// ActiveTestReq struct.
func (p *ActiveTestResp) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, nil)
}

// GetResponse implements PDU interface.
func (b *ActiveTestResp) GetResponse() codec.PDU {
	return nil
}
