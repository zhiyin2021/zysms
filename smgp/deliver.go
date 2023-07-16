package smgp

import (
	"fmt"

	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/utils"
)

type SmgpDeliveryReq struct {
	seqId      uint32
	MsgId      string // 【10字节】短消息流水号
	IsReport   byte   // 【1字节】是否为状态报告
	MsgFormat  byte   // 【1字节】短消息格式
	RecvTime   string // 【14字节】短消息定时发送时间
	SrcTermID  string // 【21字节】短信息发送方号码
	DestTermID string // 【21】短消息接收号码
	MsgLength  byte   // 【1字节】短消息长度
	MsgContent string // 【MsgLength字节】短消息内容
	// MsgBytes   []byte         // 消息内容按照Msg_Fmt编码后的数据
	//Report  *Report        // 状态报告
	Reserve string         // 【8字节】保留
	TlvList *utils.TlvList // 【TLV】可选项参数

	// 协议版本,不是报文内容，但在调用encode方法前需要设置此值
	// Version Version
}

type SmgpDeliveryRsp struct {
	seqId  uint32
	MsgId  string // 【10字节】短消息流水号
	Status Status

	// 协议版本,不是报文内容，但在调用encode方法前需要设置此值
	// Version Version
}

func (p *SmgpDeliveryReq) Pack(seqId uint32) []byte {
	pktLen := SMGP_HEADEER_LEN + 10 + 1 + 1 + 14 + 21 + 21 + 1 + uint32(p.MsgLength) + 8

	pkt := codec.NewWriter(pktLen, SMGP_DELIVER.ToInt())
	pkt.WriteU32(seqId)
	pkt.WriteU32(pktLen)
	pkt.WriteU32(SMGP_DELIVER.ToInt())
	if seqId > 0 {
		p.seqId = seqId
	}
	pkt.WriteU32(p.seqId)

	pkt.WriteStr(p.MsgId, 10)
	pkt.WriteByte(p.IsReport)
	pkt.WriteByte(p.MsgFormat)
	pkt.WriteStr(p.RecvTime, 14)
	pkt.WriteStr(p.SrcTermID, 21)
	pkt.WriteStr(p.DestTermID, 21)
	pkt.WriteByte(p.MsgLength)
	pkt.WriteStr(p.MsgContent, int(p.MsgLength))
	pkt.WriteStr(p.Reserve, 8)
	return pkt.Bytes()
}

func (p *SmgpDeliveryReq) Unpack(data []byte) error {
	var pkt = codec.NewReader(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	p.MsgId = pkt.ReadStr(10)
	p.IsReport = pkt.ReadByte()
	p.MsgFormat = pkt.ReadByte()
	p.RecvTime = pkt.ReadStr(14)
	p.SrcTermID = pkt.ReadStr(21)
	p.DestTermID = pkt.ReadStr(21)
	p.MsgLength = pkt.ReadByte()
	p.MsgContent = pkt.ReadStr(int(p.MsgLength))
	p.Reserve = pkt.ReadStr(8)
	return pkt.Err()
}
func (p *SmgpDeliveryReq) SeqId() uint32 {
	return p.seqId
}

func (p *SmgpDeliveryRsp) Pack(seqId uint32) []byte {
	pktLen := SMGP_HEADEER_LEN + 10 + 4

	pkt := codec.NewWriter(pktLen, SMGP_DELIVER_RESP.ToInt())
	pkt.WriteU32(seqId)
	p.seqId = seqId
	pkt.WriteStr(p.MsgId, 10)
	pkt.WriteU32(uint32(p.Status))
	return pkt.Bytes()
}

func (p *SmgpDeliveryRsp) Unpack(data []byte) error {
	var pkt = codec.NewReader(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	p.MsgId = pkt.ReadStr(10)
	p.Status = Status(pkt.ReadU32())
	return pkt.Err()
}

func (p *SmgpDeliveryRsp) SeqId() uint32 {
	return p.seqId
}
func (r *SmgpDeliveryRsp) String() string {
	return fmt.Sprintf("{ seq: %d, msgId: %x, status: \"%s\" }", r.SeqId, r.MsgId, r.Status)
}
