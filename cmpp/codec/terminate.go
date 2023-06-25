package codec

import (
	"github.com/zhiyin2021/zysms/proto"
)

// Packet length const for cmpp terminate request and response packets.
const (
	CmppTerminateReqLen uint32 = 12 //12d, 0xc
	CmppTerminateRspLen uint32 = 12 //12d, 0xc
)

type CmppTerminateReq struct {
	// session info
	SeqId uint32
}
type CmppTerminateRsp struct {
	// session info
	SeqId uint32
}

// Pack packs the CmppTerminateReq to bytes stream for client side.
func (p *CmppTerminateReq) Pack(seqId uint32) []byte {
	data := make([]byte, CmppTerminateReqLen)
	pkt := proto.NewPacket(data)

	// Pack header
	pkt.WriteU32(CmppTerminateReqLen)
	pkt.WriteU32(CMPP_TERMINATE.ToInt())
	pkt.WriteU32(seqId)
	p.SeqId = seqId

	return data
}

// Unpack unpack the binary byte stream to a CmppTerminateReq variable.
// After unpack, you will get all value of fields in
// CmppTerminateReq struct.
func (p *CmppTerminateReq) Unpack(data []byte) {

	pkt := proto.NewPacket(data)

	// Sequence Id
	p.SeqId = pkt.ReadU32()
}

// Pack packs the CmppTerminateRsp to bytes stream for client side.
func (p *CmppTerminateRsp) Pack(seqId uint32) []byte {
	data := make([]byte, CmppTerminateRspLen)
	pkt := proto.NewPacket(data)

	// Pack header
	pkt.WriteU32(CmppTerminateRspLen)
	pkt.WriteU32(CMPP_TERMINATE_RESP.ToInt())
	pkt.WriteU32(seqId)
	p.SeqId = seqId

	return data
}

// Unpack unpack the binary byte stream to a CmppTerminateRsp variable.
// After unpack, you will get all value of fields in
// CmppTerminateRsp struct.
func (p *CmppTerminateRsp) Unpack(data []byte) {
	pkt := proto.NewPacket(data)

	// Sequence Id
	p.SeqId = pkt.ReadU32()
}
