package cmpp

import (
	"errors"

	"github.com/zhiyin2021/zysms/event"
	"github.com/zhiyin2021/zysms/proto"
)

// Packet length const for cmpp submit request and response packets.
const (
	Cmpp2SubmitReqMaxLen uint32 = 12 + 2265  //2277d, 0x8e5
	Cmpp2SubmitRspLen    uint32 = 12 + 8 + 1 //21d, 0x15

	Cmpp3SubmitReqMaxLen uint32 = 12 + 3479  //3491d, 0xda3
	Cmpp3SubmitRspLen    uint32 = 12 + 8 + 4 //24d, 0x18
)

// Errors for result in submit resp.
var (
	ErrnoSubmitInvalidStruct         uint8 = 1
	ErrnoSubmitInvalidCommandId      uint8 = 2
	ErrnoSubmitInvalidSequence       uint8 = 3
	ErrnoSubmitInvalidMsgLength      uint8 = 4
	ErrnoSubmitInvalidFeeCode        uint8 = 5
	ErrnoSubmitExceedMaxMsgLength    uint8 = 6
	ErrnoSubmitInvalidServiceId      uint8 = 7
	ErrnoSubmitNotPassFlowControl    uint8 = 8
	ErrnoSubmitNotServeFeeTerminalId uint8 = 9
	ErrnoSubmitInvalidSrcId          uint8 = 10
	ErrnoSubmitInvalidMsgSrc         uint8 = 11
	ErrnoSubmitInvalidFeeTerminalId  uint8 = 12
	ErrnoSubmitInvalidDestTerminalId uint8 = 13

	SubmitRspResultErrMap = map[uint8]error{
		ErrnoSubmitInvalidStruct:         errSubmitInvalidStruct,
		ErrnoSubmitInvalidCommandId:      errSubmitInvalidCommandId,
		ErrnoSubmitInvalidSequence:       errSubmitInvalidSequence,
		ErrnoSubmitInvalidMsgLength:      errSubmitInvalidMsgLength,
		ErrnoSubmitInvalidFeeCode:        errSubmitInvalidFeeCode,
		ErrnoSubmitExceedMaxMsgLength:    errSubmitExceedMaxMsgLength,
		ErrnoSubmitInvalidServiceId:      errSubmitInvalidServiceId,
		ErrnoSubmitNotPassFlowControl:    errSubmitNotPassFlowControl,
		ErrnoSubmitNotServeFeeTerminalId: errSubmitNotServeFeeTerminalId,
		ErrnoSubmitInvalidSrcId:          errSubmitInvalidSrcId,
		ErrnoSubmitInvalidMsgSrc:         errSubmitInvalidMsgSrc,
		ErrnoSubmitInvalidFeeTerminalId:  errSubmitInvalidFeeTerminalId,
		ErrnoSubmitInvalidDestTerminalId: errSubmitInvalidDestTerminalId,
	}

	errSubmitInvalidStruct         = errors.New("submit response status: invalid protocol structure")
	errSubmitInvalidCommandId      = errors.New("submit response status: invalid command id")
	errSubmitInvalidSequence       = errors.New("submit response status: invalid message sequence")
	errSubmitInvalidMsgLength      = errors.New("submit response status: invalid message length")
	errSubmitInvalidFeeCode        = errors.New("submit response status: invalid fee code")
	errSubmitExceedMaxMsgLength    = errors.New("submit response status: exceed max message length")
	errSubmitInvalidServiceId      = errors.New("submit response status: invalid service id")
	errSubmitNotPassFlowControl    = errors.New("submit response status: not pass the flow control")
	errSubmitNotServeFeeTerminalId = errors.New("submit response status: feeTerminalId is not served")
	errSubmitInvalidSrcId          = errors.New("submit response status: invalid srcId")
	errSubmitInvalidMsgSrc         = errors.New("submit response status: invalid msgSrc")
	errSubmitInvalidFeeTerminalId  = errors.New("submit response status: invalid feeTerminalId")
	errSubmitInvalidDestTerminalId = errors.New("submit response status: invalid destTerminalId")
)

type Cmpp2SubmitReq struct {
	MsgId              uint64   // 信息标识，由 SP 接入的短信网关本身产生，本处填空(0)。【8字节】
	PkTotal            uint8    // 相同 Msg_Id 的信息总条数，从 1 开始。【1字节】
	PkNumber           uint8    // 相同 Msg_Id 的信息序号，从 1 开始。【1字节】
	RegisteredDelivery uint8    // 是否要求返回状态确认报告：0：不需要；1：需要。【1字节】
	MsgLevel           uint8    // 信息级别。【1字节】
	ServiceId          string   // 业务类型，是数字、字母和符号的组合。【10字节】
	FeeUserType        uint8    // 计费用户类型字段：0：对目的终端 MSISDN 计费；1：对源终端 MSISDN 计费；2：对 SP 计费；3：表示本字段无效，对谁计费参见 Fee_terminal_Id 字段。【1字节】
	FeeTerminalId      string   // 被计费用户的号码。【21字节】
	TpPid              uint8    // GSM 协议类型。详细解释请参考 GSM03.40 中的
	TpUdhi             uint8    // GSM 协议类型。详细解释请参考 GSM03.40 中的
	MsgFmt             uint8    // 信息格式：0：ASCII 码；3：短信写卡操作；4：二进制信息；8：UCS2 编码；15：含 GBK 汉字。【1字节】
	MsgSrc             string   // 信息内容来源(SP_Id)。【6字节】
	FeeType            string   // 资费类别：01：对“计费用户号码”免费；02：对“计费用户号码”按条计信息费；03：对“计费用户号码”按包月收取信息费。【2字节】
	FeeCode            string   // 资费代码（以分为单位）。【6字节】
	ValidTime          string   // 存活有效期，格式遵循 SMPP3.3 协议。【17字节】
	AtTime             string   // 定时发送时间，格式遵从 SMPP3.3 协议。【17字节】
	SrcId              string   // 源号码。【21字节】
	DestUsrTl          uint8    // 接收信息的用户数量(小于 100 个用户)。【1字节】
	DestTerminalId     []string // 接收短信的 MSISDN 号码。【21字节*DestUsr_tl】
	MsgLength          uint8    // 信息长度(Msg_Fmt 值为 0 时：< 160 个字节；其它<=140 个字节)。【1字节】
	MsgContent         string   // 信息内容。【Msg_Length 字节】
	Reserve            string   // 保留，扩展用。【8字节】
	// session info
	seqId uint32
}

type Cmpp2SubmitRsp struct {
	MsgId  uint64
	Result uint8

	// session info
	seqId uint32
}

type Cmpp3SubmitReq struct {
	MsgId              uint64
	PkTotal            uint8
	PkNumber           uint8
	RegisteredDelivery uint8
	MsgLevel           uint8
	ServiceId          string
	FeeUserType        uint8
	FeeTerminalId      string
	FeeTerminalType    uint8
	TpPid              uint8
	TpUdhi             uint8
	MsgFmt             uint8
	MsgSrc             string
	FeeType            string
	FeeCode            string
	ValidTime          string
	AtTime             string
	SrcId              string
	DestUsrTl          uint8
	DestTerminalId     []string
	DestTerminalType   uint8
	MsgLength          uint8
	MsgContent         string
	LinkId             string

	// session info
	seqId uint32
}

type Cmpp3SubmitRsp struct {
	MsgId  uint64
	Result uint32

	// session info
	seqId uint32
}

// Pack packs the Cmpp2SubmitReq to bytes stream for client side.
// Before calling Pack, you should initialize a Cmpp2SubmitReq variable
// with correct field value.
func (p *Cmpp2SubmitReq) Pack(seqId uint32) []byte {
	var pktLen uint32 = CMPP_HEADER_LEN + 117 + uint32(p.DestUsrTl)*21 + 1 + uint32(p.MsgLength) + 8
	data := make([]byte, pktLen)
	pkt := proto.NewPacket(data)

	// Pack header
	pkt.WriteU32(pktLen)
	pkt.WriteU32(CMPP_SUBMIT.ToInt())

	p.seqId = seqId

	pkt.WriteU32(p.seqId)

	// Pack Body
	pkt.WriteU64(p.MsgId)

	if p.PkTotal == 0 && p.PkNumber == 0 {
		p.PkTotal, p.PkNumber = 1, 1
	}
	pkt.WriteByte(p.PkTotal)
	pkt.WriteByte(p.PkNumber)
	pkt.WriteByte(p.RegisteredDelivery)
	pkt.WriteByte(p.MsgLevel)
	pkt.WriteStr(p.ServiceId, 10)
	pkt.WriteByte(p.FeeUserType)
	pkt.WriteStr(p.FeeTerminalId, 21)
	pkt.WriteByte(p.TpPid)
	pkt.WriteByte(p.TpUdhi)
	pkt.WriteByte(p.MsgFmt)
	pkt.WriteStr(p.MsgSrc, 6)
	pkt.WriteStr(p.FeeType, 2)
	pkt.WriteStr(p.FeeCode, 6)
	pkt.WriteStr(p.ValidTime, 17)
	pkt.WriteStr(p.AtTime, 17)
	pkt.WriteStr(p.SrcId, 21)
	pkt.WriteByte(p.DestUsrTl)

	for _, d := range p.DestTerminalId {
		pkt.WriteStr(d, 21)
	}
	pkt.WriteByte(p.MsgLength)
	pkt.WriteStr(p.MsgContent, int(p.MsgLength))
	pkt.WriteStr(p.Reserve, 8)

	return data
}

// Unpack unpack the binary byte stream to a Cmpp2SubmitReq variable.
// Usually it is used in server side. After unpack, you will get all value of fields in
// Cmpp2SubmitReq struct.
func (p *Cmpp2SubmitReq) Unpack(data []byte) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = r.(error)
		}
	}()
	pkt := proto.NewPacket(data)

	// Sequence Id
	p.seqId = pkt.ReadU32()
	p.MsgId = pkt.ReadU64()

	p.PkTotal = pkt.ReadByte()
	p.PkNumber = pkt.ReadByte()
	p.RegisteredDelivery = pkt.ReadByte()
	p.MsgLevel = pkt.ReadByte()

	serviceId := pkt.ReadStr(10)
	p.ServiceId = string(serviceId)

	p.FeeUserType = pkt.ReadByte()

	feeTerminalId := pkt.ReadStr(21)
	p.FeeTerminalId = string(feeTerminalId)

	p.TpPid = pkt.ReadByte()
	p.TpUdhi = pkt.ReadByte()
	p.MsgFmt = pkt.ReadByte()

	p.MsgSrc = pkt.ReadStr(6)

	p.FeeType = pkt.ReadStr(2)

	p.FeeCode = pkt.ReadStr(6)

	p.ValidTime = pkt.ReadStr(17)

	p.AtTime = pkt.ReadStr(17)

	p.SrcId = pkt.ReadStr(21)

	p.DestUsrTl = pkt.ReadByte()

	for i := 0; i < int(p.DestUsrTl); i++ {
		p.DestTerminalId = append(p.DestTerminalId, pkt.ReadStr(21))
	}

	p.MsgLength = pkt.ReadByte()

	p.MsgContent = pkt.ReadStr(int(p.MsgLength))

	p.Reserve = pkt.ReadStr(8)
	return nil
}
func (p *Cmpp2SubmitReq) Event() event.SmsEvent {
	return event.SmsEventSubmitReq
}
func (p *Cmpp2SubmitReq) SeqId() uint32 {
	return p.seqId
}

// Pack packs the Cmpp2SubmitRsp to bytes stream for Server side.
// Before calling Pack, you should initialize a Cmpp2SubmitRsp variable
// with correct field value.
func (p *Cmpp2SubmitRsp) Pack(seqId uint32) []byte {
	var pktLen uint32 = CMPP_HEADER_LEN + 8 + 1
	data := make([]byte, pktLen)
	pkt := proto.NewPacket(data)
	// Pack header
	pkt.WriteU32(pktLen)
	pkt.WriteU32(CMPP_SUBMIT_RESP.ToInt())

	p.seqId = seqId

	pkt.WriteU32(p.seqId)

	// Pack Body
	pkt.WriteU64(p.MsgId)
	pkt.WriteByte(p.Result)
	return data
}

// Unpack unpack the binary byte stream to a Cmpp2SubmitRsp variable.
// Usually it is used in client side. After unpack, you will get all value of fields in
// Cmpp2SubmitRsp struct.
func (p *Cmpp2SubmitRsp) Unpack(data []byte) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = r.(error)
		}
	}()
	pkt := proto.NewPacket(data)

	// Sequence Id
	p.seqId = pkt.ReadU32()
	p.MsgId = pkt.ReadU64()

	p.Result = pkt.ReadByte()
	return nil
}
func (p *Cmpp2SubmitRsp) Event() event.SmsEvent {
	return event.SmsEventSubmitRsp
}
func (p *Cmpp2SubmitRsp) SeqId() uint32 {
	return p.seqId
}

// Pack packs the Cmpp3SubmitReq to bytes stream for client side.
// Before calling Pack, you should initialize a Cmpp3SubmitReq variable
// with correct field value.
func (p *Cmpp3SubmitReq) Pack(seqId uint32) []byte {
	var pktLen uint32 = CMPP_HEADER_LEN + 129 + uint32(p.DestUsrTl)*32 + 1 + 1 + uint32(p.MsgLength) + 20
	data := make([]byte, pktLen)
	pkt := proto.NewPacket(data)
	// Pack header
	pkt.WriteU32(pktLen)
	pkt.WriteU32(CMPP_SUBMIT.ToInt())

	p.seqId = seqId

	pkt.WriteU32(p.seqId)

	// Pack Body
	pkt.WriteU64(p.MsgId)

	if p.PkTotal == 0 && p.PkNumber == 0 {
		p.PkTotal, p.PkNumber = 1, 1
	}
	pkt.WriteByte(p.PkTotal)
	pkt.WriteByte(p.PkNumber)
	pkt.WriteByte(p.RegisteredDelivery)
	pkt.WriteByte(p.MsgLevel)
	pkt.WriteStr(p.ServiceId, 10)
	pkt.WriteByte(p.FeeUserType)
	pkt.WriteStr(p.FeeTerminalId, 32)
	pkt.WriteByte(p.FeeTerminalType)
	pkt.WriteByte(p.TpPid)
	pkt.WriteByte(p.TpUdhi)
	pkt.WriteByte(p.MsgFmt)
	pkt.WriteStr(p.MsgSrc, 6)
	pkt.WriteStr(p.FeeType, 2)
	pkt.WriteStr(p.FeeCode, 6)
	pkt.WriteStr(p.ValidTime, 17)
	pkt.WriteStr(p.AtTime, 17)
	pkt.WriteStr(p.SrcId, 21)
	pkt.WriteByte(p.DestUsrTl)

	for _, d := range p.DestTerminalId {
		pkt.WriteStr(d, 32)
	}
	pkt.WriteByte(p.DestTerminalType)
	pkt.WriteByte(p.MsgLength)
	pkt.WriteStr(p.MsgContent, int(p.MsgLength))
	pkt.WriteStr(p.LinkId, 20)

	return data
}

// Unpack unpack the binary byte stream to a Cmpp3SubmitReq variable.
// Usually it is used in server side. After unpack, you will get all value of fields in
// Cmpp3SubmitReq struct.
func (p *Cmpp3SubmitReq) Unpack(data []byte) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = r.(error)
		}
	}()
	pkt := proto.NewPacket(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	p.MsgId = pkt.ReadU64()

	p.PkTotal = pkt.ReadByte()
	p.PkNumber = pkt.ReadByte()
	p.RegisteredDelivery = pkt.ReadByte()
	p.MsgLevel = pkt.ReadByte()

	p.ServiceId = pkt.ReadStr(10)

	p.FeeUserType = pkt.ReadByte()

	p.FeeTerminalId = pkt.ReadStr(32)

	p.FeeTerminalType = pkt.ReadByte()
	p.TpPid = pkt.ReadByte()
	p.TpUdhi = pkt.ReadByte()
	p.MsgFmt = pkt.ReadByte()

	p.MsgSrc = pkt.ReadStr(6)

	p.FeeType = pkt.ReadStr(2)

	p.FeeCode = pkt.ReadStr(6)

	p.ValidTime = pkt.ReadStr(17)

	p.AtTime = pkt.ReadStr(17)

	p.SrcId = pkt.ReadStr(21)

	p.DestUsrTl = pkt.ReadByte()

	for i := 0; i < int(p.DestUsrTl); i++ {
		p.DestTerminalId = append(p.DestTerminalId, pkt.ReadStr(32))
	}

	p.DestTerminalType = pkt.ReadByte()
	p.MsgLength = pkt.ReadByte()

	p.MsgContent = pkt.ReadStr(int(p.MsgLength))

	p.LinkId = pkt.ReadStr(20)
	return nil
}
func (p *Cmpp3SubmitReq) Event() event.SmsEvent {
	return event.SmsEventSubmitReq
}
func (p *Cmpp3SubmitReq) SeqId() uint32 {
	return p.seqId
}

// Pack packs the Cmpp3SubmitRsp to bytes stream for Server side.
// Before calling Pack, you should initialize a Cmpp3SubmitRsp variable
// with correct field value.
func (p *Cmpp3SubmitRsp) Pack(seqId uint32) []byte {
	var pktLen uint32 = CMPP_HEADER_LEN + 8 + 4
	data := make([]byte, pktLen)
	pkt := proto.NewPacket(data)

	// Pack header
	pkt.WriteU32(pktLen)
	pkt.WriteU32(CMPP_SUBMIT_RESP.ToInt())

	p.seqId = seqId

	pkt.WriteU32(p.seqId)

	// Pack Body
	pkt.WriteU64(p.MsgId)
	pkt.WriteU32(p.Result)

	return data
}

// Unpack unpack the binary byte stream to a Cmpp3SubmitRsp variable.
// Usually it is used in client side. After unpack, you will get all value of fields in
// Cmpp3SubmitRsp struct.
func (p *Cmpp3SubmitRsp) Unpack(data []byte) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = r.(error)
		}
	}()
	pkt := proto.NewPacket(data)

	// Sequence Id
	p.seqId = pkt.ReadU32()
	p.MsgId = pkt.ReadU64()
	p.Result = pkt.ReadU32()
	return nil
}
func (p *Cmpp3SubmitRsp) Event() event.SmsEvent {
	return event.SmsEventSubmitRsp
}

func (p *Cmpp3SubmitRsp) SeqId() uint32 {
	return p.seqId
}
