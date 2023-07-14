package cmpp

import (
	"github.com/zhiyin2021/zysms/event"
	"github.com/zhiyin2021/zysms/proto"
)

// Packet length const for cmpp terminate request and response packets.
const (
	CmppTerminateReqLen uint32 = 12 //12d, 0xc
	CmppTerminateRspLen uint32 = 12 //12d, 0xc
)

type CmppTerminateReq struct {
	// session info
	seqId uint32
}
type CmppTerminateRsp struct {
	// session info
	seqId uint32
}

// Pack packs the CmppTerminateReq to bytes stream for client side.
func (p *CmppTerminateReq) Pack(seqId uint32, sp proto.SmsProto) []byte {
	pkt := proto.NewCmppBuffer(CmppTerminateReqLen, CMPP_TERMINATE.ToInt(), seqId)
	p.seqId = seqId
	return pkt.Bytes()
}

// Unpack unpack the binary byte stream to a CmppTerminateReq variable.
// After unpack, you will get all value of fields in
// CmppTerminateReq struct.
func (p *CmppTerminateReq) Unpack(data []byte, sp proto.SmsProto) error {
	pkt := proto.NewBuffer(data)

	// Sequence Id
	p.seqId = pkt.ReadU32()
	return pkt.Err()
}
func (p *CmppTerminateReq) Event() event.SmsEvent {
	return event.SmsEventTerminateReq
}
func (p *CmppTerminateReq) SeqId() uint32 {
	return p.seqId
}

// Pack packs the CmppTerminateRsp to bytes stream for client side.
func (p *CmppTerminateRsp) Pack(seqId uint32, sp proto.SmsProto) []byte {
	pkt := proto.NewCmppBuffer(CmppTerminateRspLen, CMPP_TERMINATE_RESP.ToInt(), seqId)

	p.seqId = seqId
	return pkt.Bytes()
}

// Unpack unpack the binary byte stream to a CmppTerminateRsp variable.
// After unpack, you will get all value of fields in
// CmppTerminateRsp struct.
func (p *CmppTerminateRsp) Unpack(data []byte, sp proto.SmsProto) error {
	pkt := proto.NewBuffer(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	return pkt.Err()
}

func (p *CmppTerminateRsp) Event() event.SmsEvent {
	return event.SmsEventTerminateRsp
}
func (p *CmppTerminateRsp) SeqId() uint32 {
	return p.seqId
}
