package cmpp

import (
	"errors"

	"github.com/zhiyin2021/zysms/codec"
)

// Packet length const for cmpp deliver request and response packets.
const (
	Cmpp2DeliverReqMaxLen uint32 = 12 + 233   //245d, 0xf5
	Cmpp2DeliverRspLen    uint32 = 12 + 8 + 1 //21d, 0x15

	Cmpp3DeliverReqMaxLen uint32 = 12 + 257   //269d, 0x10d
	Cmpp3DeliverRspLen    uint32 = 12 + 8 + 4 //24d, 0x18

	ReportLen byte = 60
)

// Errors for result in deliver resp.

var (
	ErrnoDeliverInvalidStruct      uint8 = 1
	ErrnoDeliverInvalidCommandId   uint8 = 2
	ErrnoDeliverInvalidSequence    uint8 = 3
	ErrnoDeliverInvalidMsgLength   uint8 = 4
	ErrnoDeliverInvalidFeeCode     uint8 = 5
	ErrnoDeliverExceedMaxMsgLength uint8 = 6
	ErrnoDeliverInvalidServiceId   uint8 = 7
	ErrnoDeliverNotPassFlowControl uint8 = 8
	ErrnoDeliverOtherError         uint8 = 9

	DeliverRspResultErrMap = map[uint8]error{
		ErrnoDeliverInvalidStruct:      errDeliverInvalidStruct,
		ErrnoDeliverInvalidCommandId:   errDeliverInvalidCommandId,
		ErrnoDeliverInvalidSequence:    errDeliverInvalidSequence,
		ErrnoDeliverInvalidMsgLength:   errDeliverInvalidMsgLength,
		ErrnoDeliverInvalidFeeCode:     errDeliverInvalidFeeCode,
		ErrnoDeliverExceedMaxMsgLength: errDeliverExceedMaxMsgLength,
		ErrnoDeliverInvalidServiceId:   errDeliverInvalidServiceId,
		ErrnoDeliverNotPassFlowControl: errDeliverNotPassFlowControl,
		ErrnoDeliverOtherError:         errDeliverOtherError,
	}

	errDeliverInvalidStruct      = errors.New("deliver response status: invalid protocol structure")
	errDeliverInvalidCommandId   = errors.New("deliver response status: invalid command id")
	errDeliverInvalidSequence    = errors.New("deliver response status: invalid message sequence")
	errDeliverInvalidMsgLength   = errors.New("deliver response status: invalid message length")
	errDeliverInvalidFeeCode     = errors.New("deliver response status: invalid fee code")
	errDeliverExceedMaxMsgLength = errors.New("deliver response status: exceed max message length")
	errDeliverInvalidServiceId   = errors.New("deliver response status: invalid service id")
	errDeliverNotPassFlowControl = errors.New("deliver response status: not pass the flow control")
	errDeliverOtherError         = errors.New("deliver response status: other error")
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
type DeliverReport struct {
	MsgId          uint64 // 消息标识 8字节
	Stat           string // 状态 7字节
	SubmitTime     string // YYMMDDHHMM 10字节
	DoneTime       string // YYMMDDHHMM 10字节
	DestTerminalId string // 接收短信的手机号 21字节
	SmscSequence   uint32 // 短信中心的Sequence 4字节
}

func NewDeliverReq(ver Version) codec.PDU {
	c := &DeliverReq{
		base: newBase(ver, CMPP_DELIVER, 0),
	}
	return c
}
func NewDeliverResp(ver Version) codec.PDU {
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
		bw.WriteByte(p.RegisterDelivery)

		if p.RegisterDelivery == 1 && p.Report != nil {
			bw.WriteByte(ReportLen)
			bw.WriteU64(p.Report.MsgId)
			bw.WriteStr(p.Report.Stat, 7)
			bw.WriteStr(p.Report.SubmitTime, 10)
			bw.WriteStr(p.Report.DoneTime, 10)
			bw.WriteStr(p.Report.DestTerminalId, 21)
			bw.WriteU32(p.Report.SmscSequence)
		} else {
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
		p.TpPid = br.ReadByte()
		p.TpUdhi = br.ReadByte()
		p.MsgFmt = br.ReadByte()
		p.SrcTerminalId = br.ReadStr(21)
		if p.Version == V30 {
			p.SrcTerminalType = br.ReadByte()
		}
		p.RegisterDelivery = br.ReadByte()

		if p.RegisterDelivery == 1 {
			if br.ReadByte() == ReportLen {
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
			p.Result = uint32(br.ReadByte())
		}
		return br.Err()
	})
}
func (p *DeliverResp) GetResponse() codec.PDU {
	return nil
}
