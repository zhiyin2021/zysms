package cmpp

import (
	"errors"

	"github.com/zhiyin2021/zysms/proto"
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

type Cmpp2FwdReq struct {
	SourceId           string
	DestinationId      string
	NodesCount         uint8
	MsgFwdType         uint8
	MsgId              uint64
	PkTotal            uint8
	PkNumber           uint8
	RegisteredDelivery uint8
	MsgLevel           uint8
	ServiceId          string
	FeeUserType        uint8
	FeeTerminalId      string
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
	DestId             []string
	MsgLength          uint8
	MsgContent         string
	Reserve            string

	// session info
	seqId uint32
}

type Cmpp2FwdRsp struct {
	MsgId    uint64
	PkTotal  uint8
	PkNumber uint8
	Result   uint8

	// session info
	seqId uint32
}
type Cmpp3FwdReq struct {
	SourceId            string
	DestinationId       string
	NodesCount          uint8
	MsgFwdType          uint8
	MsgId               uint64
	PkTotal             uint8
	PkNumber            uint8
	RegisteredDelivery  uint8
	MsgLevel            uint8
	ServiceId           string
	FeeUserType         uint8
	FeeTerminalId       string
	FeeTerminalPseudo   string
	FeeTerminalUserType uint8
	TpPid               uint8
	TpUdhi              uint8
	MsgFmt              uint8
	MsgSrc              string
	FeeType             string
	FeeCode             string
	ValidTime           string
	AtTime              string
	SrcId               string
	SrcPseudo           string
	SrcUserType         uint8
	SrcType             uint8
	DestUsrTl           uint8
	DestId              []string
	DestPseudo          string
	DestUserType        uint8
	MsgLength           uint8
	MsgContent          string
	LinkId              string

	// session info
	seqId uint32
}

type Cmpp3FwdRsp struct {
	MsgId    uint64
	PkTotal  uint8
	PkNumber uint8
	Result   uint32

	// session info
	seqId uint32
}

// Pack packs the Cmpp2FwdReq to bytes stream for client side.
// Before calling Pack, you should initialize a Cmpp2FwdReq variable
// with correct field value.
func (p *Cmpp2FwdReq) Pack(seqId uint32) []byte {
	var pktLen uint32 = CMPP_HEADER_LEN + 131 + uint32(p.DestUsrTl)*21 + 1 + uint32(p.MsgLength) + 8
	data := make([]byte, pktLen)
	pkt := proto.NewPacket(data)

	// Pack header
	pkt.WriteU32(pktLen)
	pkt.WriteU32(CMPP_FWD.ToInt())

	p.seqId = seqId

	pkt.WriteU32(p.seqId)

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
	for _, d := range p.DestId {
		pkt.WriteStr(d, 21)
	}
	pkt.WriteByte(p.MsgLength)
	pkt.WriteStr(p.MsgContent, int(p.MsgLength))
	pkt.WriteStr(p.Reserve, 8)

	return data
}

// Unpack unpack the binary byte stream to a Cmpp2FwdReq variable.
// After unpack, you will get all value of fields in Cmpp2FwdReq struct.
func (p *Cmpp2FwdReq) Unpack(data []byte) proto.Packer {
	pkt := proto.NewPacket(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()

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
		p.DestId = append(p.DestId, pkt.ReadStr(21))
	}

	p.MsgLength = pkt.ReadByte()

	p.MsgContent = pkt.ReadStr(int(p.MsgLength))

	p.Reserve = pkt.ReadStr(8)
	return p
}
func (p *Cmpp2FwdReq) SeqId() uint32 {
	return p.seqId
}

// Pack packs the Cmpp2FwdRsp to bytes stream for server side.
// Before calling Pack, you should initialize a Cmpp2FwdRsp variable
// with correct field value.
func (p *Cmpp2FwdRsp) Pack(seqId uint32) []byte {
	data := make([]byte, Cmpp2FwdRspLen)
	pkt := proto.NewPacket(data)
	// Pack header
	pkt.WriteU32(Cmpp2FwdRspLen)
	pkt.WriteU32(CMPP_FWD_RESP.ToInt())

	p.seqId = seqId

	pkt.WriteU32(p.seqId)

	// Pack Body
	pkt.WriteU64(p.MsgId)
	pkt.WriteByte(p.PkTotal)
	pkt.WriteByte(p.PkNumber)
	pkt.WriteByte(p.Result)
	return data
}

// Unpack unpack the binary byte stream to a Cmpp2FwdRsp variable.
// After unpack, you will get all value of fields in Cmpp2FwdRsp struct.
func (p *Cmpp2FwdRsp) Unpack(data []byte) proto.Packer {
	pkt := proto.NewPacket(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	p.MsgId = pkt.ReadU64()

	p.PkTotal = pkt.ReadByte()
	p.PkNumber = pkt.ReadByte()
	p.Result = pkt.ReadByte()
	return p
}
func (p *Cmpp2FwdRsp) SeqId() uint32 {
	return p.seqId
}

// Pack packs the Cmpp3FwdReq to bytes stream for client side.
// Before calling Pack, you should initialize a Cmpp3FwdReq variable
// with correct field value.
func (p *Cmpp3FwdReq) Pack(seqId uint32) []byte {
	var pktLen uint32 = CMPP_HEADER_LEN + 198 + uint32(p.DestUsrTl)*21 + 32 + 1 + 1 + uint32(p.MsgLength) + 20
	data := make([]byte, pktLen)
	pkt := proto.NewPacket(data)

	// Pack header
	pkt.WriteU32(pktLen)
	pkt.WriteU32(CMPP_FWD.ToInt())

	p.seqId = seqId
	pkt.WriteU32(p.seqId)

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
	pkt.WriteStr(p.FeeTerminalPseudo, 32)
	pkt.WriteByte(p.FeeTerminalUserType)
	pkt.WriteByte(p.TpPid)
	pkt.WriteByte(p.TpUdhi)
	pkt.WriteByte(p.MsgFmt)
	pkt.WriteStr(p.MsgSrc, 6)
	pkt.WriteStr(p.FeeType, 2)
	pkt.WriteStr(p.FeeCode, 6)
	pkt.WriteStr(p.ValidTime, 17)
	pkt.WriteStr(p.AtTime, 17)
	pkt.WriteStr(p.SrcId, 21)
	pkt.WriteStr(p.SrcPseudo, 32)
	pkt.WriteByte(p.SrcUserType)
	pkt.WriteByte(p.SrcType)
	pkt.WriteByte(p.DestUsrTl)

	for _, d := range p.DestId {
		pkt.WriteStr(d, 21)
	}
	pkt.WriteStr(p.DestPseudo, 32)
	pkt.WriteByte(p.DestUserType)
	pkt.WriteByte(p.MsgLength)
	pkt.WriteStr(p.MsgContent, int(p.MsgLength))
	pkt.WriteStr(p.LinkId, 20)

	return data
}

// Unpack unpack the binary byte stream to a Cmpp3FwdReq variable.
// After unpack, you will get all value of fields in Cmpp3FwdReq struct.
func (p *Cmpp3FwdReq) Unpack(data []byte) proto.Packer {
	pkt := proto.NewPacket(data)
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
	p.FeeTerminalPseudo = pkt.ReadStr(32)
	p.FeeTerminalUserType = pkt.ReadByte()

	p.TpPid = pkt.ReadByte()
	p.TpUdhi = pkt.ReadByte()
	p.MsgFmt = pkt.ReadByte()

	p.MsgSrc = pkt.ReadStr(6)

	p.FeeType = pkt.ReadStr(2)

	p.FeeCode = pkt.ReadStr(6)

	p.ValidTime = pkt.ReadStr(17)

	p.AtTime = pkt.ReadStr(17)

	p.SrcId = pkt.ReadStr(21)

	p.SrcPseudo = pkt.ReadStr(32)
	p.SrcUserType = pkt.ReadByte()
	p.SrcType = pkt.ReadByte()

	p.DestUsrTl = pkt.ReadByte()
	for i := 0; i < int(p.DestUsrTl); i++ {
		p.DestId = append(p.DestId, pkt.ReadStr(21))
	}
	p.DestPseudo = pkt.ReadStr(32)
	p.DestUserType = pkt.ReadByte()

	p.MsgLength = pkt.ReadByte()

	p.MsgContent = pkt.ReadStr(int(p.MsgLength))

	p.LinkId = pkt.ReadStr(20)
	return p
}
func (p *Cmpp3FwdReq) SeqId() uint32 {
	return p.seqId
}

// Pack packs the Cmpp3FwdRsp to bytes stream for server side.
// Before calling Pack, you should initialize a Cmpp3FwdRsp variable
// with correct field value.
func (p *Cmpp3FwdRsp) Pack(seqId uint32) []byte {
	data := make([]byte, Cmpp3FwdRspLen)
	pkt := proto.NewPacket(data)
	// Pack header
	pkt.WriteU32(Cmpp3FwdRspLen)
	pkt.WriteU32(CMPP_FWD_RESP.ToInt())

	p.seqId = seqId
	pkt.WriteU32(p.seqId)

	// Pack Body
	pkt.WriteU64(p.MsgId)
	pkt.WriteByte(p.PkTotal)
	pkt.WriteByte(p.PkNumber)
	pkt.WriteU32(p.Result)
	return data
}

// Unpack unpack the binary byte stream to a Cmpp3FwdRsp variable.
// After unpack, you will get all value of fields in Cmpp3FwdRsp struct.
func (p *Cmpp3FwdRsp) Unpack(data []byte) proto.Packer {
	pkt := proto.NewPacket(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	p.MsgId = pkt.ReadU64()

	p.PkTotal = pkt.ReadByte()
	p.PkNumber = pkt.ReadByte()
	p.Result = pkt.ReadU32()
	return p
}

func (p *Cmpp3FwdRsp) SeqId() uint32 {
	return p.seqId
}
