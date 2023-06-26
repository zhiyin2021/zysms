package cmpp

import (
	"github.com/zhiyin2021/zysms/proto"
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
func (p *CmppCancelReq) Pack(seqId uint32) []byte {
	data := make([]byte, CmppCancelReqLen)
	pkt := proto.NewPacket(data)

	// Pack header
	pkt.WriteU32(CmppCancelReqLen)
	pkt.WriteU32(CMPP_CANCEL.ToInt())

	p.seqId = seqId

	pkt.WriteU32(p.seqId)

	pkt.WriteU64(p.MsgId)
	return data
}

// Unpack unpack the binary byte stream to a CmppTerminateReq variable.
// After unpack, you will get all value of fields in
// CmppTerminateReq struct.
func (p *CmppCancelReq) Unpack(data []byte) proto.Packer {
	pkt := proto.NewPacket(data)

	// Sequence Id
	p.seqId = pkt.ReadU32()
	p.MsgId = pkt.ReadU64()
	return p
}
func (p *CmppCancelReq) SeqId() uint32 {
	return p.seqId
}

// Pack packs the CmppTerminateRsp to bytes stream for client side.
func (p *CmppCancelRsp) Pack(seqId uint32) []byte {
	data := make([]byte, CmppCancelRspLen)
	pkt := proto.NewPacket(data)

	// Pack header
	pkt.WriteU32(CmppCancelRspLen)
	pkt.WriteU32(CMPP_CANCEL_RESP.ToInt())

	p.seqId = seqId
	pkt.WriteU32(p.seqId)

	// Pack body
	pkt.WriteU32(p.SuccId)
	return data
}

// Unpack unpack the binary byte stream to a CmppTerminateRsp variable.
// After unpack, you will get all value of fields in
// CmppTerminateRsp struct.
func (p *CmppCancelRsp) Unpack(data []byte) proto.Packer {
	pkt := proto.NewPacket(data)

	// Sequence Id
	p.seqId = pkt.ReadU32()
	p.SuccId = pkt.ReadU32()
	return p
}

func (p *CmppCancelRsp) SeqId() uint32 {
	return p.seqId
}