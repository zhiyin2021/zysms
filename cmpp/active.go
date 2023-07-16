package cmpp

import (
	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/event"
)

// Packet length const for cmpp active test request and response packets.
const (
	CmppActiveTestReqLen uint32 = 12     //12d, 0xc
	CmppActiveTestRspLen uint32 = 12 + 1 //13d, 0xd
)

type CmppActiveTestReq struct {
	// session info
	seqId uint32
}
type CmppActiveTestRsp struct {
	Reserved uint8
	// session info
	seqId uint32
}

// Pack packs the CmppActiveTestReq to bytes stream for client side.
func (p *CmppActiveTestReq) Pack(seqId uint32, sp codec.SmsProto) []byte {
	// buf := make([]byte, CmppActiveTestReqLen)
	pkt := codec.NewWriter(CmppActiveTestReqLen, CMPP_ACTIVE_TEST.ToInt())
	pkt.WriteU32(seqId)
	p.seqId = seqId
	return pkt.Bytes()
}

// Unpack unpack the binary byte stream to a CmppActiveTestReq variable.
// After unpack, you will get all value of fields in
// CmppActiveTestReq struct.
func (p *CmppActiveTestReq) Unpack(data []byte, sp codec.SmsProto) error {
	var r = codec.NewReader(data)
	// Sequence Id
	p.seqId = r.ReadU32()
	return r.Err()
}
func (p *CmppActiveTestReq) SeqId() uint32 {
	return p.seqId
}

func (p *CmppActiveTestReq) Event() event.SmsEvent {
	return event.SmsEventActiveTestReq
}

// Pack packs the CmppActiveTestRsp to bytes stream for client side.
func (p *CmppActiveTestRsp) Pack(seqId uint32, sp codec.SmsProto) []byte {
	pkt := codec.NewWriter(CmppActiveTestRspLen, CMPP_ACTIVE_TEST_RESP.ToInt())
	pkt.WriteU32(seqId)
	p.seqId = seqId
	pkt.WriteByte(p.Reserved)
	return pkt.Bytes()
}

// Unpack unpack the binary byte stream to a CmppActiveTestRsp variable.
// After unpack, you will get all value of fields in
// CmppActiveTestRsp struct.
func (p *CmppActiveTestRsp) Unpack(data []byte, sp codec.SmsProto) error {
	var r = codec.NewReader(data)
	// Sequence Id
	p.seqId = r.ReadU32()
	p.Reserved = r.ReadByte()
	return r.Err()
}
func (p *CmppActiveTestRsp) Event() event.SmsEvent {
	return event.SmsEventActiveTestRsp
}
func (p *CmppActiveTestRsp) SeqId() uint32 {
	return p.seqId
}
