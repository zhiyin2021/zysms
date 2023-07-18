package cmpp

import (
	"github.com/zhiyin2021/zysms/codec"
)

// Packet length const for cmpp terminate request and response packets.

type CancelReq struct {
	base
	MsgId uint64 // 8字节 信息标识
}
type CancelResp struct {
	base
	SuccId uint32
}

func NewCancelReq(ver codec.Version) codec.PDU {
	return &CancelReq{
		base: newBase(ver, CMPP_CANCEL, 0),
	}
}
func NewCancelResp(ver codec.Version) codec.PDU {
	return &CancelResp{
		base: newBase(ver, CMPP_CANCEL_RESP, 0),
	}
}

func (p *CancelReq) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteU64(p.MsgId)
	})
}

func (p *CancelReq) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.MsgId = br.ReadU64()
		return br.Err()
	})
}

func (b *CancelReq) GetResponse() codec.PDU {
	return &CancelResp{
		base: newBase(b.Version, CMPP_CANCEL_RESP, b.SequenceNumber),
	}
}

func (p *CancelResp) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteU32(p.SuccId)
	})
}

func (p *CancelResp) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.SuccId = br.ReadU32()
		return br.Err()
	})
}

func (b *CancelResp) GetResponse() codec.PDU {
	return nil
}
