package smgp

import "github.com/zhiyin2021/zysms/codec"

type ExitReq struct {
	base
}
type ExitResp struct {
	base
}

func NewExitReq(ver Version) codec.PDU {
	return &ExitReq{
		base: newBase(ver, SMGP_EXIT, 0),
	}
}
func NewExitResp(ver Version) codec.PDU {
	return &ExitResp{
		base: newBase(ver, SMGP_EXIT_RESP, 0),
	}
}

// Pack packs the ActiveTestReq to bytes stream for client side.
func (p *ExitReq) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, nil)
}

// Unpack unpack the binary byte stream to a ActiveTestReq variable.
// After unpack, you will get all value of fields in
// ActiveTestReq struct.
func (p *ExitReq) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, nil)
}

// GetResponse implements PDU interface.
func (b *ExitReq) GetResponse() codec.PDU {
	return &ExitResp{
		base: newBase(b.Version, SMGP_EXIT_RESP, b.SequenceNumber),
	}
}

// Pack packs the ActiveTestReq to bytes stream for client side.
func (p *ExitResp) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, nil)
}

// Unpack unpack the binary byte stream to a ActiveTestReq variable.
// After unpack, you will get all value of fields in
// ActiveTestReq struct.
func (p *ExitResp) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, nil)
}

// GetResponse implements PDU interface.
func (b *ExitResp) GetResponse() codec.PDU {
	return nil
}
