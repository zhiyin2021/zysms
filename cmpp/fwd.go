package cmpp

import (
	"github.com/zhiyin2021/zysms/codec"
)

type FwdReq struct {
	base
	SourceId            string // 6字节
	DestinationId       string // 6字节
	NodesCount          uint8
	MsgFwdType          uint8
	MsgId               uint64
	PkTotal             uint8
	PkNumber            uint8
	RegisteredDelivery  uint8
	MsgLevel            uint8
	ServiceId           string
	FeeUserType         uint8
	FeeTerminalId       string //被计费用户的号码(如本字节填空，则 表示本字段无效，对谁计费参见 Fee_UserType 字段。本字段与 Fee_UserType 字段互斥)  21字节
	FeeTerminalPseudo   string // 被计费用户伪码 32字节 cmpp3.0
	FeeTerminalUserType uint8  // 被计费用户号码类型,0:全球通,1神州行. 1字节 cmpp3.0
	TpPid               uint8
	TpUdhi              uint8
	MsgFmt              uint8
	MsgSrc              string
	FeeType             string
	FeeCode             string
	ValidTime           string
	AtTime              string
	SrcId               string
	SrcPseudo           string // 源号码的伪码 32字节 cmpp3.0
	SrcUserType         uint8  // 源号码的用户类型，0:全球通，1:神州行 1字节 cmpp3.0
	SrcType             uint8  // 传递给 SP 的源号码的类型，0:真实号 码;1:伪码 1字节 cmpp3.0
	DestUsrTl           uint8
	DestId              []string
	DestPseudo          string // 目的用户的伪码 32字节 cmpp3.0
	DestUserType        uint8  // 目的号码的用户类型，0:全球通，1: 神州行 1字节 cmpp3.0
	Message             codec.ShortMessage
	LinkId              string // 20字节 cmpp3.0

}

type FwdResp struct {
	base
	MsgId    uint64
	PkTotal  uint8
	PkNumber uint8
	Result   uint32 // cmpp3=4字节, cmpp2=1字节

}

func NewFwdReq(ver codec.Version) codec.PDU {
	return &FwdReq{
		base: newBase(ver, CMPP_FWD, 0),
	}
}

func NewFwdResp(ver codec.Version) codec.PDU {
	return &FwdResp{
		base: newBase(ver, CMPP_FWD_RESP, 0),
	}
}

// Pack packs the ActiveTestReq to bytes stream for client side.
func (p *FwdReq) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteStr(p.SourceId, 6)
		bw.WriteStr(p.DestinationId, 6)
		bw.WriteByte(p.NodesCount)
		bw.WriteByte(p.MsgFwdType)
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
		bw.WriteStr(p.FeeTerminalId, 21)
		if p.Version == V30 {
			bw.WriteStr(p.FeeTerminalPseudo, 32)
			bw.WriteByte(p.FeeTerminalUserType)
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
		if p.Version == V30 {
			bw.WriteStr(p.SrcPseudo, 32)
			bw.WriteByte(p.SrcUserType)
			bw.WriteByte(p.SrcType)
		}
		bw.WriteByte(p.DestUsrTl)
		for _, d := range p.DestId {
			bw.WriteStr(d, 21)
		}
		if p.Version == V30 {
			bw.WriteStr(p.DestPseudo, 32)
			bw.WriteByte(p.DestUserType)
		}
		p.Message.Marshal(bw)
		if p.Version == V30 {
			bw.WriteStr(p.LinkId, 20)
		} else {
			bw.WriteStr(p.LinkId, 8)
		}
	})
}
func (p *FwdReq) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.SourceId = br.ReadStr(6)
		p.DestinationId = br.ReadStr(6)
		p.NodesCount = br.ReadU8()
		p.MsgFwdType = br.ReadU8()
		p.MsgId = br.ReadU64()
		p.PkTotal = br.ReadU8()
		p.PkNumber = br.ReadU8()
		p.RegisteredDelivery = br.ReadU8()
		p.MsgLevel = br.ReadU8()
		p.ServiceId = br.ReadStr(10)
		p.FeeUserType = br.ReadU8()
		p.FeeTerminalId = br.ReadStr(21)
		if p.Version == V30 {
			p.FeeTerminalPseudo = br.ReadStr(32)
			p.FeeTerminalUserType = br.ReadU8()
		}
		p.TpPid = br.ReadU8()
		p.TpUdhi = br.ReadU8()
		p.MsgFmt = br.ReadU8()
		p.MsgSrc = br.ReadStr(6)
		p.FeeType = br.ReadStr(2)
		p.FeeCode = br.ReadStr(6)
		p.ValidTime = br.ReadStr(17)
		p.AtTime = br.ReadStr(17)
		p.SrcId = br.ReadStr(21)
		if p.Version == V30 {
			p.SrcPseudo = br.ReadStr(32)
			p.SrcUserType = br.ReadU8()
			p.SrcType = br.ReadU8()
		}
		p.DestUsrTl = br.ReadU8()
		for i := 0; i < int(p.DestUsrTl); i++ {
			p.DestId = append(p.DestId, br.ReadStr(21))
		}
		if p.Version == V30 {
			p.DestPseudo = br.ReadStr(32)
			p.DestUserType = br.ReadU8()
		}
		p.Message.Unmarshal(br, p.TpUdhi != 0, p.MsgFmt)
		if p.Version == V30 {
			p.LinkId = br.ReadStr(20)
		} else {
			p.LinkId = br.ReadStr(8)
		}
		return br.Err()
	})
}

// Pack packs the ActiveTestReq to bytes stream for client side.
func (p *FwdResp) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteU64(p.MsgId)
		bw.WriteByte(p.PkTotal)
		bw.WriteByte(p.PkNumber)
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
func (p *FwdResp) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.MsgId = br.ReadU64()
		p.PkTotal = br.ReadU8()
		p.PkNumber = br.ReadU8()
		if p.Version == V30 {
			p.Result = br.ReadU32()
		} else {
			p.Result = uint32(br.ReadU8())
		}
		return br.Err()
	})
}

func (p *FwdReq) GetResponse() codec.PDU {
	return &FwdResp{
		base: newBase(p.Version, CMPP_FWD_RESP, p.SequenceNumber),
	}
}

func (p *FwdResp) GetResponse() codec.PDU {
	return nil
}
