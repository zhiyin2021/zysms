package sgip

import "github.com/zhiyin2021/zysms/codec"

type UnbindReq struct {
	base
}
type UnbindResp struct {
	base
}

func NewUnbindReq(ver codec.Version, nodeId uint32) codec.PDU {
	return &UnbindReq{
		base: newBase(ver, SGIP_UNBIND, [3]uint32{nodeId, 0, 0}),
	}
}
func NewUnbindResp(ver codec.Version, nodeId uint32) codec.PDU {
	return &UnbindResp{
		base: newBase(ver, SGIP_UNBIND_RESP, [3]uint32{nodeId, 0, 0}),
	}
}
func (p *UnbindReq) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, nil)
}

func (p *UnbindReq) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, nil)
}

func (b *UnbindReq) GetResponse() codec.PDU {
	return &UnbindResp{
		base: newBase(b.Version, SGIP_UNBIND_RESP, b.SequenceNumber),
	}
}

func (p *UnbindResp) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, nil)
}

func (p *UnbindResp) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, nil)
}

func (b *UnbindResp) GetResponse() codec.PDU {
	return nil
}
