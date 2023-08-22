package smgp

import (
	"github.com/zhiyin2021/zysms/codec"
)

type SubmitReq struct {
	base
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

	Message codec.ShortMessage // 消息内容按照Msg_Fmt编码后的数据
	Reserve string             // 【8字节】保留

}

type SubmitResp struct {
	base
	MsgId  string // 【10字节】短消息流水号
	Status Status
}

func NewSubmitReq(ver codec.Version) codec.PDU {
	return &SubmitReq{
		base: newBase(ver, SMGP_SUBMIT, 0),
	}
}
func NewSubmitResp(ver codec.Version) codec.PDU {
	return &SubmitResp{
		base: newBase(ver, SMGP_SUBMIT_RESP, 0),
	}
}

// Pack packs the ActiveTestReq to bytes stream for client side.
func (p *SubmitReq) Marshal(w *codec.BytesWriter) {
	if p.Message.IsLongMessage() {
		p.RegisterOptionalParam(codec.NewTlv(codec.TagTPUdhi, []byte{0x01}))
	}
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteByte(p.SubType)
		bw.WriteByte(p.NeedReport)
		bw.WriteByte(p.Priority)
		bw.WriteStr(p.ServiceID, 10)
		bw.WriteStr(p.FeeType, 2)
		bw.WriteStr(p.FeeCode, 6)
		bw.WriteStr(p.FixedFee, 6)
		bw.WriteByte(p.MsgFormat)
		bw.WriteStr(p.ValidTime, 17)
		bw.WriteStr(p.AtTime, 17)
		bw.WriteStr(p.SrcTermID, 21)
		bw.WriteStr(p.ChargeTermID, 21)
		bw.WriteByte(p.DestTermIDCount)
		for _, tid := range p.DestTermID {
			bw.WriteStr(tid, 21)
		}
		p.Message.Marshal(bw)
		bw.WriteStr(p.Reserve, 8)
	})
}

// Unpack unpack the binary byte stream to a ActiveTestReq variable.
// After unpack, you will get all value of fields in
// ActiveTestReq struct.
func (p *SubmitReq) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.SubType = br.ReadByte()
		p.NeedReport = br.ReadByte()
		p.Priority = br.ReadByte()
		p.ServiceID = br.ReadStr(10)
		p.FeeType = br.ReadStr(2)
		p.FeeCode = br.ReadStr(6)
		p.FixedFee = br.ReadStr(6)
		p.MsgFormat = br.ReadByte()
		p.ValidTime = br.ReadStr(17)
		p.AtTime = br.ReadStr(17)
		p.SrcTermID = br.ReadStr(21)
		p.ChargeTermID = br.ReadStr(21)
		p.DestTermIDCount = br.ReadByte()
		for i := byte(0); i < p.DestTermIDCount; i++ {
			p.DestTermID = append(p.DestTermID, br.ReadStr(21))
		}
		// 05   00   03   00   04   01   #   长短信设置
		// 0002   0001   40     #   TP_udhi
		// 0009   0001   04     #   pkTotal
		// 000a   0001   01     #   pkNumber

		p.Message.Unmarshal(br, p.TpUdhi(), p.MsgFormat)
		p.Reserve = br.ReadStr(8)
		return br.Err()
	})
}
func (p *SubmitReq) TpUdhi() bool {
	udhi := byte(0)
	if tag, ok := p.OptionalParameters[codec.TagTPUdhi]; ok && len(tag.Data) > 0 {
		udhi = tag.Data[0]
	}
	return udhi == 1
}

// GetResponse implements PDU interface.
func (b *SubmitReq) GetResponse() codec.PDU {
	return &SubmitResp{
		base: newBase(b.Version, SMGP_SUBMIT_RESP, b.SequenceNumber),
	}
}

// Pack packs the ActiveTestReq to bytes stream for client side.
func (p *SubmitResp) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteStr(p.MsgId, 10)
		bw.WriteU32(uint32(p.Status))
	})
}

// Unpack unpack the binary byte stream to a ActiveTestReq variable.
// After unpack, you will get all value of fields in
// ActiveTestReq struct.
func (p *SubmitResp) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.MsgId = br.ReadStr(10)
		p.Status = Status(br.ReadU32())
		return br.Err()
	})
}

// GetResponse implements PDU interface.
func (b *SubmitResp) GetResponse() codec.PDU {
	return nil
}
