package cmpp

import (
	"errors"

	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/event"
)

const (
	Cmpp2FwdReqMaxLen uint32 = 12 + 2379          //2277d, 0x957
	Cmpp2FwdRspLen    uint32 = 12 + 8 + 1 + 1 + 1 //23d, 0x17

	Cmpp3FwdReqMaxLen uint32 = 12 + 2491          //2503d, 0x9c7
	Cmpp3FwdRspLen    uint32 = 12 + 8 + 1 + 1 + 4 //26d, 0x1a
)

// Errors for result in fwd resp.
var (
	ErrnoFwdInvalidStruct      uint8 = 1
	ErrnoFwdInvalidCommandId   uint8 = 2
	ErrnoFwdInvalidSequence    uint8 = 3
	ErrnoFwdInvalidMsgLength   uint8 = 4
	ErrnoFwdInvalidFeeCode     uint8 = 5
	ErrnoFwdExceedMaxMsgLength uint8 = 6
	ErrnoFwdInvalidServiceId   uint8 = 7
	ErrnoFwdNotPassFlowControl uint8 = 8
	ErrnoFwdNoPrivilege        uint8 = 9

	FwdRspResultErrMap = map[uint8]error{
		ErrnoFwdInvalidStruct:      errFwdInvalidStruct,
		ErrnoFwdInvalidCommandId:   errFwdInvalidCommandId,
		ErrnoFwdInvalidSequence:    errFwdInvalidSequence,
		ErrnoFwdInvalidMsgLength:   errFwdInvalidMsgLength,
		ErrnoFwdInvalidFeeCode:     errFwdInvalidFeeCode,
		ErrnoFwdExceedMaxMsgLength: errFwdExceedMaxMsgLength,
		ErrnoFwdInvalidServiceId:   errFwdInvalidServiceId,
		ErrnoFwdNotPassFlowControl: errFwdNotPassFlowControl,
		ErrnoFwdNoPrivilege:        errFwdNoPrivilege,
	}

	errFwdInvalidStruct      = errors.New("fwd response status: invalid protocol structure")
	errFwdInvalidCommandId   = errors.New("fwd response status: invalid command id")
	errFwdInvalidSequence    = errors.New("fwd response status: invalid message sequence")
	errFwdInvalidMsgLength   = errors.New("fwd response status: invalid message length")
	errFwdInvalidFeeCode     = errors.New("fwd response status: invalid fee code")
	errFwdExceedMaxMsgLength = errors.New("fwd response status: exceed max message length")
	errFwdInvalidServiceId   = errors.New("fwd response status: invalid service id")
	errFwdNotPassFlowControl = errors.New("fwd response status: not pass the flow control")
	errFwdNoPrivilege        = errors.New("fwd response status: msg has no fwd privilege")
)

// type Cmpp2FwdReq struct {
// 	SourceId           string
// 	DestinationId      string
// 	NodesCount         uint8
// 	MsgFwdType         uint8
// 	MsgId              uint64
// 	PkTotal            uint8
// 	PkNumber           uint8
// 	RegisteredDelivery uint8
// 	MsgLevel           uint8
// 	ServiceId          string
// 	FeeUserType        uint8
// 	FeeTerminalId      string
// 	TpPid              uint8
// 	TpUdhi             uint8
// 	MsgFmt             uint8
// 	MsgSrc             string
// 	FeeType            string
// 	FeeCode            string
// 	ValidTime          string
// 	AtTime             string
// 	SrcId              string
// 	DestUsrTl          uint8
// 	DestId             []string
// 	MsgLength          uint8
// 	MsgContent         string
// 	Reserve            string

// 	// session info
// 	seqId uint32
// }

type CmppFwdReq struct {
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
	MsgLength           uint8
	MsgContent          string
	LinkId              string // 20字节 cmpp3.0

	// session info
	seqId uint32
}

type CmppFwdRsp struct {
	MsgId    uint64
	PkTotal  uint8
	PkNumber uint8
	Result   uint32 // cmpp3=4字节, cmpp2=1字节

	// session info
	seqId uint32
}

// Pack packs the Cmpp3FwdReq to bytes stream for client side.
// Before calling Pack, you should initialize a Cmpp3FwdReq variable
// with correct field value.
func (p *CmppFwdReq) Pack(seqId uint32, sp codec.SmsProto) []byte {
	var pktLen uint32 = CMPP_HEADER_LEN + 131 + uint32(p.DestUsrTl)*21 + 1 + uint32(p.MsgLength) + 8
	if sp == codec.CMPP30 {
		pktLen = CMPP_HEADER_LEN + 198 + uint32(p.DestUsrTl)*21 + 32 + 1 + 1 + uint32(p.MsgLength) + 20
	}
	pkt := codec.NewWriter(pktLen, CMPP_FWD.ToInt())
	pkt.WriteU32(seqId)

	p.seqId = seqId

	// Pack Body
	pkt.WriteStr(p.SourceId, 6)
	pkt.WriteStr(p.DestinationId, 6)
	pkt.WriteByte(p.NodesCount)
	pkt.WriteByte(p.MsgFwdType)
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
	if sp == codec.CMPP30 {
		pkt.WriteStr(p.FeeTerminalPseudo, 32)
		pkt.WriteByte(p.FeeTerminalUserType)
	}
	pkt.WriteByte(p.TpPid)
	pkt.WriteByte(p.TpUdhi)
	pkt.WriteByte(p.MsgFmt)
	pkt.WriteStr(p.MsgSrc, 6)
	pkt.WriteStr(p.FeeType, 2)
	pkt.WriteStr(p.FeeCode, 6)
	pkt.WriteStr(p.ValidTime, 17)
	pkt.WriteStr(p.AtTime, 17)
	pkt.WriteStr(p.SrcId, 21)
	if sp == codec.CMPP30 {
		pkt.WriteStr(p.SrcPseudo, 32)
		pkt.WriteByte(p.SrcUserType)
		pkt.WriteByte(p.SrcType)
	}
	pkt.WriteByte(p.DestUsrTl)
	for _, d := range p.DestId {
		pkt.WriteStr(d, 21)
	}
	if sp == codec.CMPP30 {
		pkt.WriteStr(p.DestPseudo, 32)
		pkt.WriteByte(p.DestUserType)
	}
	pkt.WriteByte(p.MsgLength)
	pkt.WriteStr(p.MsgContent, int(p.MsgLength))
	if sp == codec.CMPP30 {
		pkt.WriteStr(p.LinkId, 20)
	} else {
		pkt.WriteStr(p.LinkId, 8)
	}

	return pkt.Bytes()
}

// Unpack unpack the binary byte stream to a Cmpp3FwdReq variable.
// After unpack, you will get all value of fields in Cmpp3FwdReq struct.
func (p *CmppFwdReq) Unpack(data []byte, sp codec.SmsProto) (e error) {
	pkt := codec.NewReader(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	// Body
	p.SourceId = pkt.ReadStr(6)
	p.DestinationId = pkt.ReadStr(6)
	p.NodesCount = pkt.ReadByte()
	p.MsgFwdType = pkt.ReadByte()

	p.MsgId = pkt.ReadU64()
	p.PkTotal = pkt.ReadByte()
	p.PkNumber = pkt.ReadByte()
	p.RegisteredDelivery = pkt.ReadByte()
	p.MsgLevel = pkt.ReadByte()

	p.ServiceId = pkt.ReadStr(10)

	p.FeeUserType = pkt.ReadByte()

	p.FeeTerminalId = pkt.ReadStr(21)
	if sp == codec.CMPP30 {
		p.FeeTerminalPseudo = pkt.ReadStr(32)
		p.FeeTerminalUserType = pkt.ReadByte()
	}

	p.TpPid = pkt.ReadByte()
	p.TpUdhi = pkt.ReadByte()
	p.MsgFmt = pkt.ReadByte()

	p.MsgSrc = pkt.ReadStr(6)

	p.FeeType = pkt.ReadStr(2)

	p.FeeCode = pkt.ReadStr(6)

	p.ValidTime = pkt.ReadStr(17)

	p.AtTime = pkt.ReadStr(17)

	p.SrcId = pkt.ReadStr(21)
	if sp == codec.CMPP30 {
		p.SrcPseudo = pkt.ReadStr(32)
		p.SrcUserType = pkt.ReadByte()
		p.SrcType = pkt.ReadByte()
	}

	p.DestUsrTl = pkt.ReadByte()
	for i := 0; i < int(p.DestUsrTl); i++ {
		p.DestId = append(p.DestId, pkt.ReadStr(21))
	}
	if sp == codec.CMPP30 {
		p.DestPseudo = pkt.ReadStr(32)
		p.DestUserType = pkt.ReadByte()
	}
	p.MsgLength = pkt.ReadByte()

	p.MsgContent = pkt.ReadStr(int(p.MsgLength))
	if sp == codec.CMPP30 {
		p.LinkId = pkt.ReadStr(20)
	} else {
		p.LinkId = pkt.ReadStr(8)
	}
	return nil
}
func (p *CmppFwdReq) Event() event.SmsEvent {
	return event.SmsEventFwdReq
}

func (p *CmppFwdReq) SeqId() uint32 {
	return p.seqId
}

// Pack packs the Cmpp3FwdRsp to bytes stream for server side.
// Before calling Pack, you should initialize a Cmpp3FwdRsp variable
// with correct field value.
func (p *CmppFwdRsp) Pack(seqId uint32, sp codec.SmsProto) []byte {
	rspLen := Cmpp2FwdRspLen
	if sp == codec.CMPP30 {
		rspLen = Cmpp3FwdRspLen
	}
	pkt := codec.NewWriter(rspLen, CMPP_FWD_RESP.ToInt())
	pkt.WriteU32(seqId)
	p.seqId = seqId

	// Pack Body
	pkt.WriteU64(p.MsgId)
	pkt.WriteByte(p.PkTotal)
	pkt.WriteByte(p.PkNumber)
	if sp == codec.CMPP30 {
		pkt.WriteU32(p.Result)
	} else {
		pkt.WriteByte(byte(p.Result))
	}
	return pkt.Bytes()
}

// Unpack unpack the binary byte stream to a Cmpp3FwdRsp variable.
// After unpack, you will get all value of fields in Cmpp3FwdRsp struct.
func (p *CmppFwdRsp) Unpack(data []byte, sp codec.SmsProto) error {
	pkt := codec.NewReader(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	p.MsgId = pkt.ReadU64()

	p.PkTotal = pkt.ReadByte()
	p.PkNumber = pkt.ReadByte()
	if sp == codec.CMPP30 {
		p.Result = pkt.ReadU32()
	} else {
		p.Result = uint32(pkt.ReadByte())
	}
	return pkt.Err()
}
func (p *CmppFwdRsp) Event() event.SmsEvent {
	return event.SmsEventFwdRsp
}

func (p *CmppFwdRsp) SeqId() uint32 {
	return p.seqId
}
