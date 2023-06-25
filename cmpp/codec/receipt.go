package codec

import (
	"github.com/zhiyin2021/zysms/proto"
)

// Packet length const for cmpp receipt packet.
const (
	CmppReceiptPktLen uint32 = 60 //60d, 0x3c
)

type CmppReceiptPkt struct {
	MsgId          uint64 // 信息标识，SP提交短信(CMPP_SUBMIT)操作时，与SP相连的ISMG产生的 Msg_Id。【8字节】
	Stat           string // 发送短信的应答结果。【7字节】
	SubmitTime     string // yyMMddHHmm 【10字节】
	DoneTime       string // yyMMddHHmm 【10字节】
	DestTerminalId string // SP 发送 CMPP_SUBMIT 消息的目标终端  【21字节】
	SmscSequence   uint32 // 取自SMSC发送状态报告的消息体中的消息标识。【4字节】
}

// Pack packs the CmppReceiptPkt to bytes stream for client side.
func (p *CmppReceiptPkt) Pack() []byte {
	data := make([]byte, CmppReceiptPktLen)
	pkt := proto.NewPacket(data)

	pkt.WriteU64(p.MsgId)
	pkt.WriteStr(p.Stat, 7)
	pkt.WriteStr(p.SubmitTime, 10)
	pkt.WriteStr(p.DoneTime, 10)
	pkt.WriteStr(p.DestTerminalId, 21)
	pkt.WriteU32(p.SmscSequence)

	return data
}

// Unpack unpack the binary byte stream to a CmppReceiptPkt variable.
// After unpack, you will get all value of fields in
// CmppReceiptPkt struct.
func (p *CmppReceiptPkt) Unpack(data []byte) {
	pkt := proto.NewPacket(data)

	p.MsgId = pkt.ReadU64()
	p.Stat = pkt.ReadStr(7)
	p.SubmitTime = pkt.ReadStr(10)
	p.DoneTime = pkt.ReadStr(10)
	p.DestTerminalId = pkt.ReadStr(21)
	p.SmscSequence = pkt.ReadU32()

}
