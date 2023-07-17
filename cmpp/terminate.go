package cmpp

import (
	"github.com/zhiyin2021/zysms/codec"
)

type TerminateReq struct {
	base
}
type TerminateResp struct {
	base
}

func NewTerminateReq(ver Version) codec.PDU {
	c := &TerminateReq{
		base: newBase(ver, CMPP_TERMINATE, 0),
	}
	return c
}
func NewTerminateResp(ver Version) codec.PDU {
	c := &TerminateResp{
		base: newBase(ver, CMPP_TERMINATE_RESP, 0),
	}
	return c
}

// Pack packs the ActiveTestReq to bytes stream for client side.
func (p *TerminateReq) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, nil)
}

// Unpack unpack the binary byte stream to a ActiveTestReq variable.
// After unpack, you will get all value of fields in
// ActiveTestReq struct.
func (p *TerminateReq) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, nil)
}
func (p *TerminateReq) GetResponse() codec.PDU {
	return &TerminateResp{
		base: newBase(p.Version, CMPP_TERMINATE_RESP, 0),
	}
}

// Pack packs the ActiveTestReq to bytes stream for client side.
func (p *TerminateResp) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, nil)
}

// Unpack unpack the binary byte stream to a ActiveTestReq variable.
// After unpack, you will get all value of fields in
// ActiveTestReq struct.
func (p *TerminateResp) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, nil)
}
func (p *TerminateResp) GetResponse() codec.PDU {
	return nil
}
