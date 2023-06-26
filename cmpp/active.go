package cmpp

import (
	"github.com/zhiyin2021/zysms/proto"
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
func (p *CmppActiveTestReq) Pack(seqId uint32) []byte {
	buf := make([]byte, CmppActiveTestReqLen)
	pkt := proto.NewPacket(buf)
	// Pack header
	pkt.WriteU32(CmppActiveTestReqLen)
	pkt.WriteU32(CMPP_ACTIVE_TEST.ToInt())
	p.seqId = seqId
	pkt.WriteU32(p.seqId)
	return buf
}

// Unpack unpack the binary byte stream to a CmppActiveTestReq variable.
// After unpack, you will get all value of fields in
// CmppActiveTestReq struct.
func (p *CmppActiveTestReq) Unpack(data []byte) proto.Packer {
	var r = proto.NewPacket(data)
	// Sequence Id
	p.seqId = r.ReadU32()
	return p
}
func (p *CmppActiveTestReq) SeqId() uint32 {
	return p.seqId
}

// Pack packs the CmppActiveTestRsp to bytes stream for client side.
func (p *CmppActiveTestRsp) Pack(seqId uint32) []byte {
	buf := make([]byte, CmppActiveTestRspLen)
	pkt := proto.NewPacket(buf)
	// Pack header
	pkt.WriteU32(CmppActiveTestRspLen)
	pkt.WriteU32(CMPP_ACTIVE_TEST_RESP.ToInt())

	p.seqId = seqId
	pkt.WriteU32(p.seqId)
	pkt.WriteByte(p.Reserved)
	return buf
}

// Unpack unpack the binary byte stream to a CmppActiveTestRsp variable.
// After unpack, you will get all value of fields in
// CmppActiveTestRsp struct.
func (p *CmppActiveTestRsp) Unpack(data []byte) proto.Packer {
	var r = proto.NewPacket(data)
	// Sequence Id
	p.seqId = r.ReadU32()
	p.Reserved = r.ReadByte()
	return p
}
func (p *CmppActiveTestRsp) SeqId() uint32 {
	return p.seqId
}
