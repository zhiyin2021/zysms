package smgp

import (
	"fmt"

	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/utils"
)

type SmgpSubmitReq struct {
	seqId           uint32
	SubType         byte     // 【1字节】短消息类型
	NeedReport      byte     // 【1字节】SP是否要求返回状态报告
	Priority        byte     // 【1字节】短消息发送优先级
	ServiceID       string   // 【10字节】业务代码
	FeeType         string   // 【2字节】收费类型
	FeeCode         string   // 【6字节】资费代码
	FixedFee        string   // 【6字节】包月费/封顶费
	MsgFormat       byte     // 【1字节】短消息格式
	ValidTime       string   // 【17字节】短消息有效时间
	AtTime          string   // 【17字节】短消息定时发送时间
	SrcTermID       string   // 【21字节】短信息发送方号码
	ChargeTermID    string   // 【21字节】计费用户号码
	DestTermIDCount byte     // 【1字节】短消息接收号码总数
	DestTermID      []string // 【21*DestTermCount字节】短消息接收号码
	MsgLength       byte     // 【1字节】短消息长度
	MsgContent      string   // 消息内容按照Msg_Fmt编码后的数据
	Reserve         string   // 【8字节】保留

	TlvList *utils.TlvList // 【TLV】可选项参数
}

type SmgpSubmitRsp struct {
	seqId  uint32
	MsgId  string // 【10字节】短消息流水号
	Status Status
}

const MtBaseLen = 126

func (p *SmgpSubmitReq) Pack(seqId uint32) []byte {
	pktLen := SMGP_HEADEER_LEN + 117 + uint32(p.DestTermIDCount)*21 + 1 + uint32(p.MsgLength) + 8
	pkt := codec.NewWriter(pktLen, SMGP_SUBMIT.ToInt())
	pkt.WriteU32(seqId)

	p.seqId = seqId

	pkt.WriteByte(p.SubType)
	pkt.WriteByte(p.NeedReport)
	pkt.WriteByte(p.Priority)
	pkt.WriteStr(p.ServiceID, 10)
	pkt.WriteStr(p.FeeType, 2)
	pkt.WriteStr(p.FeeCode, 6)
	pkt.WriteStr(p.FixedFee, 6)
	pkt.WriteByte(p.MsgFormat)
	pkt.WriteStr(p.ValidTime, 17)
	pkt.WriteStr(p.AtTime, 17)
	pkt.WriteStr(p.SrcTermID, 21)
	pkt.WriteStr(p.ChargeTermID, 21)
	pkt.WriteByte(p.DestTermIDCount)
	for _, tid := range p.DestTermID {
		pkt.WriteStr(tid, 21)
	}
	pkt.WriteByte(p.MsgLength)
	pkt.WriteStr(p.MsgContent, int(p.MsgLength))
	pkt.WriteStr(p.Reserve, 8)
	return pkt.Bytes()
}

func (p *SmgpSubmitReq) Unpack(data []byte) error {
	pkt := codec.NewReader(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	// Body: Source_Addr
	p.SubType = pkt.ReadByte()
	p.NeedReport = pkt.ReadByte()
	p.Priority = pkt.ReadByte()
	p.ServiceID = pkt.ReadStr(10)
	p.FeeType = pkt.ReadStr(2)
	p.FeeCode = pkt.ReadStr(6)
	p.FixedFee = pkt.ReadStr(6)
	p.MsgFormat = pkt.ReadByte()
	p.ValidTime = pkt.ReadStr(17)
	p.AtTime = pkt.ReadStr(17)
	p.SrcTermID = pkt.ReadStr(21)
	p.ChargeTermID = pkt.ReadStr(21)
	p.DestTermIDCount = pkt.ReadByte()
	for i := byte(0); i < p.DestTermIDCount; i++ {
		p.DestTermID = append(p.DestTermID, pkt.ReadStr(21))
	}
	p.MsgLength = pkt.ReadByte()
	p.MsgContent = pkt.ReadStr(int(p.MsgLength))
	p.Reserve = pkt.ReadStr(8)
	return pkt.Err()
}
func (p *SmgpSubmitReq) SeqId() uint32 {
	return p.seqId
}

func (p *SmgpSubmitReq) String() string {
	bts := p.MsgContent
	if p.MsgLength > 6 {
		bts = p.MsgContent[:6]
	}
	return fmt.Sprintf("{ seq: %d, subType: %v, NeedReport: %v, LruPriority: %v, ServiceID: %v, "+
		"feeType: %v, feeCode: %v, fixedFee: %v, msgFormat: %v, validTime: %v, AtTime: %v, SrcTermID: %v, "+
		"chargeTermID: %v, destTermIDCount: %v, destTermID: %v, msgLength: %v, msgContent: %#x..., "+
		"reserve: %v}",
		p.seqId, p.SubType, p.NeedReport, p.Priority, p.ServiceID,
		p.FeeType, p.FeeCode, p.FixedFee, p.MsgFormat, p.ValidTime, p.AtTime, p.SrcTermID,
		p.ChargeTermID, p.DestTermIDCount, p.DestTermID, p.MsgLength, bts,
		p.Reserve)
}

func (p *SmgpSubmitRsp) Pack(seqId uint32) []byte {
	pktLen := SMGP_HEADEER_LEN + 10 + 4
	pkt := codec.NewWriter(pktLen, SMGP_SUBMIT_RESP.ToInt())
	pkt.WriteU32(seqId)
	p.seqId = seqId
	pkt.WriteStr(p.MsgId, 10)
	pkt.WriteU32(uint32(p.Status))
	return pkt.Bytes()
}

func (p *SmgpSubmitRsp) Unpack(data []byte) error {
	pkt := codec.NewReader(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	// Body: Source_Addr
	p.MsgId = pkt.ReadStr(10)
	p.Status = Status(pkt.ReadU32())
	return pkt.Err()
}

func (p *SmgpSubmitRsp) SeqId() uint32 {
	return p.seqId
}
