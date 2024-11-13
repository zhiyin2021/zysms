package cmpp

import (
	"fmt"

	"github.com/zhiyin2021/zysms/codec"
)

// Packet length const for cmpp deliver request and response packets.
const (
	ReportLen byte = 60
)

type DeliverReq struct {
	base
	MsgId            uint64 // 消息标识
	DestId           string // 目的号码 21
	ServiceId        string // 业务类型 10
	TpPid            uint8  // GSM协议类型。详细是解释请参考GSM03.40中的9.2.3.9 【1字节】
	TpUdhi           uint8  // GSM协议类型。详细是解释请参考GSM03.40中的9.2.3.9 【1字节】
	MsgFmt           uint8  // 信息格式 【1字节】
	SrcTerminalId    string // 源终端MSISDN号码（状态报告时填为CMPP_SUBMIT消息的目的终端号码）
	SrcTerminalType  uint8  // 源终端号码类型，0：真实号码；1：伪码  cmpp3.0新增项 【1字节】
	RegisterDelivery uint8  // 是否为状态报告
	// MsgLength        uint8
	// MsgContent       string
	Message codec.ShortMessage
	LinkId  string // cmpp3.0 = 20字节 点播业务, cmpp2.0 = 8字节 保留项
	Report  *DeliverReport
}
type DeliverResp struct {
	base
	MsgId  uint64
	Result uint32 // cmpp3.0 = 4字节, cmpp2.0 = 1字节
}

func NewDeliverReq(ver codec.Version) codec.PDU {
	c := &DeliverReq{
		base: newBase(ver, CMPP_DELIVER, 0),
	}
	return c
}
func NewDeliverResp(ver codec.Version) codec.PDU {
	c := &DeliverResp{
		base: newBase(ver, CMPP_DELIVER_RESP, 0),
	}
	return c
}

// Pack packs the Cmpp3DeliverReq to bytes stream for client side.
func (p *DeliverReq) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteU64(p.MsgId)
		bw.WriteStr(p.DestId, 21)
		bw.WriteStr(p.ServiceId, 10)
		bw.WriteByte(p.TpPid)
		bw.WriteByte(p.TpUdhi)
		bw.WriteByte(p.MsgFmt)
		bw.WriteStr(p.SrcTerminalId, 21)
		if p.Version == V30 {
			bw.WriteByte(p.SrcTerminalType)
		}

		if p.Report != nil {
			bw.WriteByte(1) // p.RegisterDelivery = 1
			bw.WriteByte(ReportLen)
			bw.WriteU64(p.Report.MsgId)
			bw.WriteStr(p.Report.Stat, 7)
			bw.WriteStr(p.Report.SubmitTime, 10)
			bw.WriteStr(p.Report.DoneTime, 10)
			bw.WriteStr(p.Report.DestTerminalId, 21)
			bw.WriteU32(p.Report.SmscSequence)
		} else {
			bw.WriteByte(0) // p.RegisterDelivery = 0
			p.Message.Marshal(bw)
		}
		if p.Version == V30 {
			bw.WriteStr(p.LinkId, 20)
		} else {
			// cmpp2 写入reserved 保留字段8字节
			bw.WriteStr(p.LinkId, 8)
		}

	})
}

// Unpack unpack the binary byte stream to a Cmpp3DeliverReq variable.
// After unpack, you will get all value of fields in
// Cmpp3DeliverReq struct.
func (p *DeliverReq) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.MsgId = br.ReadU64()
		p.DestId = br.ReadStr(21)
		p.ServiceId = br.ReadStr(10)
		p.TpPid = br.ReadU8()
		p.TpUdhi = br.ReadU8()
		p.MsgFmt = br.ReadU8()
		p.SrcTerminalId = br.ReadStr(21)
		if p.Version == V30 {
			p.SrcTerminalType = br.ReadU8()
		}
		p.RegisterDelivery = br.ReadU8()

		if p.RegisterDelivery == 1 {
			if br.ReadU8() == ReportLen {
				p.Report = &DeliverReport{
					MsgId:          br.ReadU64(),
					Stat:           br.ReadStr(7),
					SubmitTime:     br.ReadStr(10),
					DoneTime:       br.ReadStr(10),
					DestTerminalId: br.ReadStr(21),
					SmscSequence:   br.ReadU32(),
				}
			}
		} else {
			p.Message.Unmarshal(br, p.TpUdhi == 1, p.MsgFmt)
		}
		if p.Version == V30 {
			p.LinkId = br.ReadStr(20)
		} else {
			// cmpp2 读取reserved 保留字段8字节
			p.LinkId = br.ReadStr(8)
		}
		return br.Err()
	})
}
func (req DeliverReq) String() string {
	return fmt.Sprintf("deliverReq%s msgId:%d,src:%v,dst:%v,fmt:%d,msg:%x,rep:%s,opts:%s", req.Header, req.MsgId, req.SrcTerminalId, req.DestId, req.Message.DataCoding(), req.Message.GetMessageData(), req.Report, req.OptionalParameters)
}

func (p *DeliverReq) GetResponse() codec.PDU {
	return &DeliverResp{
		base: newBase(p.Version, CMPP_DELIVER_RESP, p.SequenceNumber),
	}
}

// Pack packs the Cmpp3DeliverRsp to bytes stream for client side.
func (p *DeliverResp) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteU64(p.MsgId)
		if p.Version == V30 {
			bw.WriteU32(p.Result)
		} else {
			bw.WriteByte(byte(p.Result))
		}
	})
}

// Unpack unpack the binary byte stream to a Cmpp3DeliverRsp variable.
// After unpack, you will get all value of fields in
// Cmpp3DeliverRsp struct.
func (p *DeliverResp) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.MsgId = br.ReadU64()
		if p.Version == V30 {
			p.Result = br.ReadU32()
		} else {
			p.Result = uint32(br.ReadU8())
		}
		return br.Err()
	})
}
func (p *DeliverResp) GetResponse() codec.PDU {
	return nil
}
func (resp DeliverResp) String() string {
	return fmt.Sprintf("deliverResp:%s msgId:%d,stat:%d,opts:%s", resp.Header, resp.MsgId, resp.Result, resp.OptionalParameters)
}

type DeliverReport struct {
	MsgId          uint64 // 消息标识 8字节
	Stat           string // 状态 7字节
	SubmitTime     string // YYMMDDHHMM 10字节
	DoneTime       string // YYMMDDHHMM 10字节
	DestTerminalId string // 接收短信的手机号 21字节
	SmscSequence   uint32 // 短信中心的Sequence 4字节
}

func (r *DeliverReport) String() string {
	return fmt.Sprintf("id:%d stat:%s submit date:%s done date:%s dest:%s smsc:%d ", r.MsgId, r.Stat, r.SubmitTime, r.DoneTime, r.DestTerminalId, r.SmscSequence)

}
