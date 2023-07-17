package cmpp

import (
	"github.com/zhiyin2021/zysms/codec"
)

/*
6位协议头格式：05 00 03 XX MM NN
byte 1 : 05, 表示剩余协议头的长度
byte 2 : 00, 这个值在GSM 03.40规范9.2.3.24.1中规定，表示随后的这批超长短信的标识位长度为1（格式中的XX值）。
byte 3 : 03, 这个值表示剩下短信标识的长度
byte 4 : XX，这批短信的唯一标志(被拆分的多条短信,此值必需一致)，事实上，SME(手机或者SP)把消息合并完之后，就重新记录，所以这个标志是否唯
一并不是很 重要。
byte 5 : MM, 这批短信的数量。如果一个超长短信总共5条，这里的值就是5。
byte 6 : NN, 这批短信的数量。如果当前短信是这批短信中的第一条的值是1，第二条的值是2。
例如：05 00 03 39 02 01
*/

type SubmitReq struct {
	base
	MsgId              uint64
	PkTotal            uint8
	PkNumber           uint8
	RegisteredDelivery uint8
	MsgLevel           uint8
	ServiceId          string
	FeeUserType        uint8
	FeeTerminalId      string // cmpp3.0=32字节, cmpp2.0=21字节
	FeeTerminalType    uint8  // 被计费用户的号码类型，0：真实号码；1：伪码。【1字节】cmpp3.0
	TpPid              uint8
	TpUdhi             uint8 // 0：消息内容体不带协议头；1：消息内容体带协议头。【1字节】
	MsgFmt             uint8
	MsgSrc             string
	FeeType            string
	FeeCode            string
	ValidTime          string
	AtTime             string
	SrcId              string
	DestUsrTl          uint8
	DestTerminalId     []string // 接收短信的 MSISDN 号码。cmpp3.0 = 32字节*DestUsr_tl, cmpp2.0 = 21*DestUsr_tl
	DestTerminalType   uint8    // 接收短信的用户的号码类型，0：真实号码；1：伪码。【1字节】 cmpp3.0
	// MsgLength          uint8
	Message codec.ShortMessage //[]byte // 信息内容
	LinkId  string             // 点播业务使用 LinkID,非点播业务的MT流程不使用该字段  cmpp3.0 = 20字节, cmpp2.0 = 8字节
}

type SubmitResp struct {
	base
	MsgId  uint64
	Result uint32 // 3.0 = 4字节, 2.0 = 1字节
}

func NewSubmitReq(ver Version) codec.PDU {
	c := &SubmitReq{
		base: newBase(ver, CMPP_SUBMIT, 0),
	}
	return c
}
func NewSubmitResp(ver Version) codec.PDU {
	c := &SubmitResp{
		base: newBase(ver, CMPP_SUBMIT_RESP, 0),
	}
	return c
}

func (p *SubmitReq) numLen() int {
	if p.Version == V30 {
		return 32
	}
	return 21
}

// Pack packs the ActiveTestReq to bytes stream for client side.
func (p *SubmitReq) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteU64(p.MsgId)
		if p.PkTotal == 0 && p.PkNumber == 0 {
			p.PkTotal, p.PkNumber = 1, 1
		}
		bw.WriteByte(p.PkTotal)
		bw.WriteByte(p.PkNumber)
		bw.WriteByte(p.RegisteredDelivery)
		bw.WriteByte(p.MsgLevel)
		bw.WriteStr(p.ServiceId, 10)
		bw.WriteByte(p.FeeUserType)
		bw.WriteStr(p.FeeTerminalId, p.numLen()) // numLen => cmpp3.0=32字节, cmpp2.0=21字节
		if p.Version == V30 {
			bw.WriteByte(p.FeeTerminalType)
		}
		bw.WriteByte(p.TpPid)
		bw.WriteByte(p.TpUdhi)
		bw.WriteByte(p.MsgFmt)
		bw.WriteStr(p.MsgSrc, 6)
		bw.WriteStr(p.FeeType, 2)
		bw.WriteStr(p.FeeCode, 6)
		bw.WriteStr(p.ValidTime, 17)
		bw.WriteStr(p.AtTime, 17)
		bw.WriteStr(p.SrcId, 21)
		p.DestUsrTl = uint8(len(p.DestTerminalId))
		bw.WriteByte(p.DestUsrTl)

		for _, d := range p.DestTerminalId {
			bw.WriteStr(d, p.numLen()) // numLen => cmpp3.0=32字节, cmpp2.0=21字节
		}
		if p.Version == V30 {
			bw.WriteByte(p.DestTerminalType)
		}
		// bw.WriteByte(p.MsgLength)
		// bw.WriteBytes(p.MsgContent)
		p.Message.Marshal(bw)
		if p.Version == V30 {
			bw.WriteStr(p.LinkId, 20)
		} else {
			bw.WriteStr(p.LinkId, 8)
		}
	})
}

// Unpack unpack the binary byte stream to a ActiveTestReq variable.
// After unpack, you will get all value of fields in
// ActiveTestReq struct.
func (p *SubmitReq) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.MsgId = br.ReadU64()
		p.PkTotal = br.ReadByte()
		p.PkNumber = br.ReadByte()
		p.RegisteredDelivery = br.ReadByte()
		p.MsgLevel = br.ReadByte()
		p.ServiceId = br.ReadStr(10)
		p.FeeUserType = br.ReadByte()
		p.FeeTerminalId = br.ReadStr(p.numLen())
		if p.Version == V30 {
			p.FeeTerminalType = br.ReadByte()
		}
		p.TpPid = br.ReadByte()
		p.TpUdhi = br.ReadByte()
		p.MsgFmt = br.ReadByte()
		p.MsgSrc = br.ReadStr(6)
		p.FeeType = br.ReadStr(2)
		p.FeeCode = br.ReadStr(6)
		p.ValidTime = br.ReadStr(17)
		p.AtTime = br.ReadStr(17)
		p.SrcId = br.ReadStr(21)
		p.DestUsrTl = br.ReadByte()
		for i := 0; i < int(p.DestUsrTl); i++ {
			p.DestTerminalId = append(p.DestTerminalId, br.ReadStr(p.numLen()))
		}
		if p.Version == V30 {
			p.DestTerminalType = br.ReadByte()
		}
		p.Message.Unmarshal(br, p.TpUdhi == 1, p.MsgFmt)
		if p.Version == V30 {
			p.LinkId = br.ReadStr(20)
		} else {
			p.LinkId = br.ReadStr(8)
		}
		return br.Err()
	})
}

func (p *SubmitReq) GetResponse() codec.PDU {
	return &SubmitResp{
		base: newBase(p.Version, CMPP_SUBMIT_RESP, p.SequenceNumber),
	}
}

// Pack packs the ActiveTestReq to bytes stream for client side.
func (p *SubmitResp) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteU64(p.MsgId)
		if p.Version == V30 {
			bw.WriteU32(p.Result)
		} else {
			bw.WriteByte(byte(p.Result))
		}
	})
}

// Unpack unpack the binary byte stream to a ActiveTestReq variable.
// After unpack, you will get all value of fields in
// ActiveTestReq struct.
func (p *SubmitResp) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.MsgId = br.ReadU64()
		if p.Version == V30 {
			p.Result = br.ReadU32()
		} else {
			p.Result = uint32(br.ReadByte())
		}
		return br.Err()
	})
}

func (p *SubmitResp) GetResponse() codec.PDU {
	return nil
}
