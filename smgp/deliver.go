package smgp

import (
	"github.com/zhiyin2021/zysms/codec"
)

type DeliveryReq struct {
	base
	MsgId      string // 【10字节】短消息流水号
	IsReport   byte   // 【1字节】是否为状态报告
	MsgFormat  byte   // 【1字节】短消息格式
	RecvTime   string // 【14字节】短消息定时发送时间
	SrcTermID  string // 【21字节】短信息发送方号码
	DestTermID string // 【21】短消息接收号码

	Message codec.ShortMessage // 【MsgLength字节】短消息内容
	// MsgBytes   []byte         // 消息内容按照Msg_Fmt编码后的数据
	//Report  *Report        // 状态报告
	Reserve string // 【8字节】保留

	// 协议版本,不是报文内容，但在调用encode方法前需要设置此值
	// Version Version
}

type DeliveryResp struct {
	base
	MsgId  string // 【10字节】短消息流水号
	Status Status

	// 协议版本,不是报文内容，但在调用encode方法前需要设置此值
	// Version Version
}

func NewDeliveryReq(ver Version) codec.PDU {
	return &DeliveryReq{
		base: newBase(ver, SMGP_DELIVER, 0),
	}
}
func NewDeliveryResp(ver Version) codec.PDU {
	return &DeliveryResp{
		base: newBase(ver, SMGP_DELIVER_RESP, 0),
	}
}

// Pack packs the ActiveTestReq to bytes stream for client side.
func (p *DeliveryReq) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteStr(p.MsgId, 10)
		bw.WriteByte(p.IsReport)
		bw.WriteByte(p.MsgFormat)
		bw.WriteStr(p.RecvTime, 14)
		bw.WriteStr(p.SrcTermID, 21)
		bw.WriteStr(p.DestTermID, 21)
		p.Message.Marshal(bw)

		bw.WriteStr(p.Reserve, 8)
	})
}

// Unpack unpack the binary byte stream to a ActiveTestReq variable.
// After unpack, you will get all value of fields in
// ActiveTestReq struct.
func (p *DeliveryReq) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.MsgId = br.ReadStr(10)
		p.IsReport = br.ReadByte()
		p.MsgFormat = br.ReadByte()
		p.RecvTime = br.ReadStr(14)
		p.SrcTermID = br.ReadStr(21)
		p.DestTermID = br.ReadStr(21)
		udhi := byte(0)
		if tag, ok := p.OptionalParameters[codec.TagTPUdhi]; ok && len(tag.Data) > 0 {
			udhi = tag.Data[0]
		}
		p.Message.Unmarshal(br, udhi == 1, p.MsgFormat)
		p.Reserve = br.ReadStr(8)
		return br.Err()
	})
}

// GetResponse implements PDU interface.
func (b *DeliveryReq) GetResponse() codec.PDU {
	return &DeliveryResp{
		base: newBase(b.Version, SMGP_DELIVER_RESP, b.SequenceNumber),
	}
}

// Pack packs the ActiveTestReq to bytes stream for client side.
func (p *DeliveryResp) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteStr(p.MsgId, 10)
		bw.WriteU32(uint32(p.Status))
	})
}

// Unpack unpack the binary byte stream to a ActiveTestReq variable.
// After unpack, you will get all value of fields in
// ActiveTestReq struct.
func (p *DeliveryResp) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.MsgId = br.ReadStr(10)
		p.Status = Status(br.ReadU32())
		return br.Err()
	})
}

// GetResponse implements PDU interface.
func (b *DeliveryResp) GetResponse() codec.PDU {
	return nil
}
