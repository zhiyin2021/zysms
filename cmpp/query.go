package cmpp

import (
	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/event"
)

// Packet length const for cmpp terminate request and response packets.
const (
	CmppQueryReqLen uint32 = 12 + 27 //12d, 0xc
	CmppQueryRspLen uint32 = 12 + 51 //12d, 0xc
)

type CmppQueryReq struct {
	// session info
	seqId     uint32
	Time      string //	8字节 YYYYMMDD
	QueryType byte   //	1字节 0：总数查询；1：按业务类型查询
	QueryCode string // 10字节  查询码 当 QueryType 为 0 时，此项无效;当 QueryType 为 1 时，此项填写业务类 型 Service_Id.
	Reserve   string // 8字节 保留
}
type CmppQueryRsp struct {
	// session info
	seqId     uint32
	Time      string //	8字节 YYYYMMDD
	QueryType byte   //	1字节 0：总数查询；1：按业务类型查询
	QueryCode string // 10字节  查询码 当 QueryType 为 0 时，此项无效;当 QueryType 为 1 时，此项填写业务类 型 Service_Id.

	MtTlMsg uint32 // 4字节 从SP接收信息总数
	MtTlUsr uint32 // 4字节 从SP接收用户总数
	MtScs   uint32 // 4字节 成功转发数量
	MtWt    uint32 // 4字节 待转发数量
	MtFl    uint32 // 4字节 转发失败数量
	MoScs   uint32 // 4字节 向SP成功送达数量
	MoWt    uint32 // 4字节 向SP待送达数量
	MoFl    uint32 // 4字节 向SP送达失败数量

}

// Pack packs the CmppTerminateReq to bytes stream for client side.
func (p *CmppQueryReq) Pack(seqId uint32, sp codec.SmsProto) []byte {
	pkt := codec.NewWriter(CmppQueryReqLen, CMPP_QUERY.ToInt())
	pkt.WriteU32(seqId)

	p.seqId = seqId

	// Pack body
	pkt.WriteStr(p.Time, 8)
	pkt.WriteByte(p.QueryType)
	pkt.WriteStr(p.QueryCode, 10)
	pkt.WriteStr(p.Reserve, 8)
	return pkt.Bytes()
}

// Unpack unpack the binary byte stream to a CmppTerminateReq variable.
// After unpack, you will get all value of fields in
// CmppTerminateReq struct.
func (p *CmppQueryReq) Unpack(data []byte, sp codec.SmsProto) error {
	pkt := codec.NewReader(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	p.Time = pkt.ReadStr(8)
	p.QueryType = pkt.ReadByte()
	p.QueryCode = pkt.ReadStr(10)
	p.Reserve = pkt.ReadStr(8)
	return pkt.Err()
}
func (p *CmppQueryReq) Event() event.SmsEvent {
	return event.SmsEventQueryReq
}

func (p *CmppQueryReq) SeqId() uint32 {
	return p.seqId
}

// Pack packs the CmppTerminateRsp to bytes stream for client side.
func (p *CmppQueryRsp) Pack(seqId uint32, sp codec.SmsProto) []byte {
	pkt := codec.NewWriter(CmppQueryRspLen, CMPP_QUERY_RESP.ToInt())
	pkt.WriteU32(seqId)

	p.seqId = seqId

	// Pack body
	pkt.WriteStr(p.Time, 8)
	pkt.WriteByte(p.QueryType)
	pkt.WriteStr(p.QueryCode, 10)
	pkt.WriteU32(p.MtTlMsg)
	pkt.WriteU32(p.MtTlUsr)
	pkt.WriteU32(p.MtScs)
	pkt.WriteU32(p.MtWt)
	pkt.WriteU32(p.MtFl)
	pkt.WriteU32(p.MoScs)
	pkt.WriteU32(p.MoWt)
	pkt.WriteU32(p.MoFl)
	return pkt.Bytes()
}

// Unpack unpack the binary byte stream to a CmppTerminateRsp variable.
// After unpack, you will get all value of fields in
// CmppTerminateRsp struct.
func (p *CmppQueryRsp) Unpack(data []byte, sp codec.SmsProto) error {

	pkt := codec.NewReader(data)

	// Sequence Id
	p.seqId = pkt.ReadU32()
	p.Time = pkt.ReadStr(8)
	p.QueryType = pkt.ReadByte()
	p.QueryCode = pkt.ReadStr(10)
	p.MtTlMsg = pkt.ReadU32()
	p.MtTlUsr = pkt.ReadU32()
	p.MtScs = pkt.ReadU32()
	p.MtWt = pkt.ReadU32()
	p.MtFl = pkt.ReadU32()
	p.MoScs = pkt.ReadU32()
	p.MoWt = pkt.ReadU32()
	p.MoFl = pkt.ReadU32()
	return pkt.Err()
}

func (p *CmppQueryRsp) Event() event.SmsEvent {
	return event.SmsEventQueryRsp
}

func (p *CmppQueryRsp) SeqId() uint32 {
	return p.seqId
}
