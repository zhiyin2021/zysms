package cmpp

import (
	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/event"
)

// Packet length const for cmpp terminate request and response packets.
const (
	CmppCancelReqLen uint32 = 12 + 8 //12d, 0xc
	CmppCancelRspLen uint32 = 12 + 4 //12d, 0xc
)

type CmppCancelReq struct {
	seqId uint32
	MsgId uint64 // 8字节 信息标识
}
type CmppCancelRsp struct {
	// session info
	seqId  uint32
	SuccId uint32
}

// Pack packs the CmppTerminateReq to bytes stream for client side.
func (p *CmppCancelReq) Pack(seqId uint32, sp codec.SmsProto) []byte {
	pkt := codec.NewWriter(CmppCancelReqLen, CMPP_CANCEL.ToInt())
	pkt.WriteU32(seqId)
	p.seqId = seqId
	pkt.WriteU64(p.MsgId)
	return pkt.Bytes()
}

// Unpack unpack the binary byte stream to a CmppTerminateReq variable.
// After unpack, you will get all value of fields in
// CmppTerminateReq struct.
func (p *CmppCancelReq) Unpack(data []byte, sp codec.SmsProto) (e error) {
	pkt := codec.NewReader(data)
	// pkt := codec.NewPacket(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	p.MsgId = pkt.ReadU64()
	return pkt.Err()
}
func (p *CmppCancelReq) SeqId() uint32 {
	return p.seqId
}
func (p *CmppCancelReq) Event() event.SmsEvent {
	return event.SmsEventCancelReq
}

// Pack packs the CmppTerminateRsp to bytes stream for client side.
func (p *CmppCancelRsp) Pack(seqId uint32, sp codec.SmsProto) []byte {
	pkt := codec.NewWriter(CmppCancelRspLen, CMPP_CANCEL_RESP.ToInt())
	pkt.WriteU32(seqId)
	p.seqId = seqId
	// Pack body
	pkt.WriteU32(p.SuccId)
	return pkt.Bytes()
}

// Unpack unpack the binary byte stream to a CmppTerminateRsp variable.
// After unpack, you will get all value of fields in
// CmppTerminateRsp struct.
func (p *CmppCancelRsp) Unpack(data []byte, sp codec.SmsProto) error {
	pkt := codec.NewReader(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	p.SuccId = pkt.ReadU32()

	return pkt.Err()
}
func (p *CmppCancelRsp) Event() event.SmsEvent {
	return event.SmsEventCancelRsp
}

func (p *CmppCancelRsp) SeqId() uint32 {
	return p.seqId
}
