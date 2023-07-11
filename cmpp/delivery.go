package cmpp

import (
	"errors"

	"github.com/zhiyin2021/zysms/event"
	"github.com/zhiyin2021/zysms/proto"
)

// Packet length const for cmpp deliver request and response packets.
const (
	Cmpp2DeliverReqMaxLen uint32 = 12 + 233   //245d, 0xf5
	Cmpp2DeliverRspLen    uint32 = 12 + 8 + 1 //21d, 0x15

	Cmpp3DeliverReqMaxLen uint32 = 12 + 257   //269d, 0x10d
	Cmpp3DeliverRspLen    uint32 = 12 + 8 + 4 //24d, 0x18
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

// type Cmpp2DeliverReq struct {
// 	MsgId            uint64 // 信息标识，由 SP 接入的短信网关本身产生，本处填空(0)。【8字节】
// 	DestId           string // 目的号码 21
// 	ServiceId        string // 业务类型 10
// 	TpPid            uint8  // GSM协议类型。详细是解释请参考GSM03.40中的9.2.3.9 【1字节】
// 	TpUdhi           uint8  // GSM协议类型。详细是解释请参考GSM03.40中的9.2.3.9 【1字节】
// 	MsgFmt           uint8  // 信息格式 【1字节】
// 	SrcTerminalId    string // 源终端MSISDN号码 【21字节】
// 	RegisterDelivery uint8  // 是否为状态报告
// 	MsgLength        uint8  // 信息长度
// 	MsgContent       string // 信息内容
// 	Reserve          string // 保留

// 	Report *CmppDeliverReport
// 	//session info
// 	seqId uint32 // sequence id
// }

// type Cmpp2DeliverRsp struct {
// 	MsgId  uint64
// 	Result uint8

//		//session info
//		seqId uint32
//	}
type CmppDeliverReq struct {
	MsgId            uint64 // 消息标识
	DestId           string // 目的号码 21
	ServiceId        string // 业务类型 10
	TpPid            uint8  // GSM协议类型。详细是解释请参考GSM03.40中的9.2.3.9 【1字节】
	TpUdhi           uint8  // GSM协议类型。详细是解释请参考GSM03.40中的9.2.3.9 【1字节】
	MsgFmt           uint8  // 信息格式 【1字节】
	SrcTerminalId    string // 源终端MSISDN号码（状态报告时填为CMPP_SUBMIT消息的目的终端号码）
	SrcTerminalType  uint8  // 源终端号码类型，0：真实号码；1：伪码  cmpp3.0新增项 【1字节】
	RegisterDelivery uint8  // 是否为状态报告
	MsgLength        uint8
	MsgContent       string
	LinkId           string // cmpp3.0 = 20字节 点播业务, cmpp2.0 = 8字节 保留项
	Report           *CmppDeliverReport
	//session info
	seqId uint32
}
type CmppDeliverRsp struct {
	MsgId  uint64
	Result uint32 // cmpp3.0 = 4字节, cmpp2.0 = 1字节

	//session info
	seqId uint32
}
type CmppDeliverReport struct {
	MsgId          uint64 // 消息标识 8字节
	Stat           string // 状态 7字节
	SubmitTime     string // YYMMDDHHMM 10字节
	DoneTime       string // YYMMDDHHMM 10字节
	DestTerminalId string // 接收短信的手机号 21字节
	SmscSequence   uint32 // 短信中心的Sequence 4字节
}

// Pack packs the Cmpp3DeliverReq to bytes stream for client side.
func (p *CmppDeliverReq) Pack(seqId uint32, sp proto.SmsProto) []byte {
	var pktLen uint32 = CMPP_HEADER_LEN + 65 + uint32(p.MsgLength) + 8
	if sp == proto.CMPP30 {
		pktLen = CMPP_HEADER_LEN + 77 + uint32(p.MsgLength) + 20
	}
	data := make([]byte, pktLen)
	pkt := proto.NewPacket(data)

	// Pack header
	pkt.WriteU32(pktLen)
	pkt.WriteU32(CMPP_DELIVER.ToInt())

	p.seqId = seqId

	pkt.WriteU32(p.seqId)

	// Pack Body
	pkt.WriteU64(p.MsgId)
	pkt.WriteStr(p.DestId, 21)
	pkt.WriteStr(p.ServiceId, 10)
	pkt.WriteByte(p.TpPid)
	pkt.WriteByte(p.TpUdhi)
	pkt.WriteByte(p.MsgFmt)
	pkt.WriteStr(p.SrcTerminalId, 32)
	if sp == proto.CMPP30 {
		pkt.WriteByte(p.SrcTerminalType)
	}
	pkt.WriteByte(p.RegisterDelivery)

	if p.RegisterDelivery == 1 && p.Report != nil {
		pkt.WriteByte(60)
		pkt.WriteU64(p.Report.MsgId)
		pkt.WriteStr(p.Report.Stat, 7)
		pkt.WriteStr(p.Report.SubmitTime, 10)
		pkt.WriteStr(p.Report.DoneTime, 10)
		pkt.WriteStr(p.Report.DestTerminalId, 21)
		pkt.WriteU32(p.Report.SmscSequence)
	} else {
		pkt.WriteByte(p.MsgLength)
		pkt.WriteStr(p.MsgContent, int(p.MsgLength))
	}
	if sp == proto.CMPP30 {
		pkt.WriteStr(p.LinkId, 20)
	} else {
		// cmpp2 写入reserved 保留字段8字节
		pkt.WriteStr(p.LinkId, 8)
	}

	return data
}

// Unpack unpack the binary byte stream to a Cmpp3DeliverReq variable.
// After unpack, you will get all value of fields in
// Cmpp3DeliverReq struct.
func (p *CmppDeliverReq) Unpack(data []byte, sp proto.SmsProto) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = r.(error)
		}
	}()
	pkt := proto.NewPacket(data)

	// Sequence Id
	p.seqId = pkt.ReadU32()

	// Body
	p.MsgId = pkt.ReadU64()

	p.DestId = pkt.ReadStr(21)

	p.ServiceId = pkt.ReadStr(10)
	p.TpPid = pkt.ReadByte()
	p.TpUdhi = pkt.ReadByte()
	p.MsgFmt = pkt.ReadByte()

	p.SrcTerminalId = pkt.ReadStr(32)
	if sp == proto.CMPP30 {
		p.SrcTerminalType = pkt.ReadByte()
	}
	p.RegisterDelivery = pkt.ReadByte()
	p.MsgLength = pkt.ReadByte()

	if p.RegisterDelivery == 1 {
		p.Report = &CmppDeliverReport{
			MsgId:          pkt.ReadU64(),
			Stat:           pkt.ReadStr(7),
			SubmitTime:     pkt.ReadStr(10),
			DoneTime:       pkt.ReadStr(10),
			DestTerminalId: pkt.ReadStr(21),
			SmscSequence:   pkt.ReadU32(),
		}
	} else {
		// 0：ASCII 码；3：短信写卡操作；4：二进制信息；8：UCS2 编码；15：含 GBK 汉字。【1字节】
		if p.MsgFmt == 8 {
			p.MsgContent = pkt.ReadUCS2(int(p.MsgLength))
		} else {
			p.MsgContent = pkt.ReadStr(int(p.MsgLength))
		}
	}
	if sp == proto.CMPP30 {
		p.LinkId = pkt.ReadStr(20)
	} else {
		// cmpp2 读取reserved 保留字段8字节
		p.LinkId = pkt.ReadStr(8)
	}
	return nil
}
func (p *CmppDeliverReq) Event() event.SmsEvent {
	return event.SmsEventDeliverReq
}

func (p *CmppDeliverReq) SeqId() uint32 {
	return p.seqId
}

// Pack packs the Cmpp3DeliverRsp to bytes stream for client side.
func (p *CmppDeliverRsp) Pack(seqId uint32, sp proto.SmsProto) []byte {
	rspLen := Cmpp2DeliverRspLen
	if sp == proto.CMPP30 {
		rspLen = Cmpp3DeliverRspLen
	}
	data := make([]byte, rspLen)
	pkt := proto.NewPacket(data)

	// Pack header
	pkt.WriteU32(rspLen)
	pkt.WriteU32(CMPP_DELIVER_RESP.ToInt())

	p.seqId = seqId

	pkt.WriteU32(p.seqId)

	// Pack Body
	pkt.WriteU64(p.MsgId)
	if sp == proto.CMPP30 {
		pkt.WriteU32(p.Result)
	} else {
		pkt.WriteByte(byte(p.Result))
	}
	return data
}

// Unpack unpack the binary byte stream to a Cmpp3DeliverRsp variable.
// After unpack, you will get all value of fields in
// Cmpp3DeliverRsp struct.
func (p *CmppDeliverRsp) Unpack(data []byte, sp proto.SmsProto) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = r.(error)
		}
	}()
	pkt := proto.NewPacket(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	p.MsgId = pkt.ReadU64()
	if sp == proto.CMPP30 {
		p.Result = pkt.ReadU32()
	} else {
		p.Result = uint32(pkt.ReadByte())
	}
	return nil
}
func (p *CmppDeliverRsp) Event() event.SmsEvent {
	return event.SmsEventDeliverRsp
}

func (p *CmppDeliverRsp) SeqId() uint32 {
	return p.seqId
}
