package cmpp

import (
	"github.com/zhiyin2021/zysms/codec"
)

// Packet length const for cmpp active test request and response packets.
const (
	ActiveTestReqLen  uint32 = 12     //12d, 0xc
	ActiveTestRespLen uint32 = 12 + 1 //13d, 0xd
)

type ActiveTestReq struct {
	base
}
type ActiveTestResp struct {
	base
	Reserved uint8
}

func NewActiveTestReq(ver codec.Version) codec.PDU {
	return &ActiveTestReq{
		base: newBase(ver, CMPP_ACTIVE_TEST, 0),
	}
}

func NewActiveTestResp(ver codec.Version) codec.PDU {
	return &ActiveTestResp{
		base: newBase(ver, CMPP_ACTIVE_TEST_RESP, 0),
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
		base: newBase(b.Version, CMPP_ACTIVE_TEST_RESP, b.SequenceNumber),
	}
}

// Pack packs the ActiveTestResp to bytes stream for client side.
func (p *ActiveTestResp) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteByte(p.Reserved)
	})
}

// Unpack unpack the binary byte stream to a ActiveTestResp variable.
// After unpack, you will get all value of fields in
// ActiveTestResp struct.
func (p *ActiveTestResp) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.Reserved = br.ReadU8()
		return br.Err()
	})
}

// GetResponse implements PDU interface.
func (b *ActiveTestResp) GetResponse() codec.PDU {
	return nil
}
