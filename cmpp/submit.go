package cmpp

import (
	"errors"

	"github.com/sirupsen/logrus"
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

type CmppSubmitReq struct {
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
	TpUdhi             uint8
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
	MsgLength          uint8
	MsgContent         string
	LinkId             string // 点播业务使用 LinkID,非点播业务的MT流程不使用该字段  cmpp3.0 = 20字节, cmpp2.0 = 8字节

	// session info
	seqId uint32
}

type CmppSubmitRsp struct {
	MsgId  uint64
	Result uint32 // 3.0 = 4字节, 2.0 = 1字节

	// session info
	seqId uint32
}

func (p *CmppSubmitReq) Pack(seqId uint32, sp proto.SmsProto) []byte {
	var pktLen uint32 = CMPP_HEADER_LEN + 129 + uint32(len(p.DestTerminalId)*32) + 1 + 1 + uint32(p.MsgLength) + 20
	numLen := 32
	if sp == proto.CMPP2 {
		pktLen = CMPP_HEADER_LEN + 117 + uint32(len(p.DestTerminalId))*21 + 1 + uint32(p.MsgLength) + 8
		numLen = 21
	}
	logrus.Infof("submit.pack %d,%d", pktLen, numLen)
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
	pkt.WriteStr(p.FeeTerminalId, numLen) // numLen => cmpp3.0=32字节, cmpp2.0=21字节
	if sp == proto.CMPP3 {
		pkt.WriteByte(p.FeeTerminalType)
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
	p.DestUsrTl = uint8(len(p.DestTerminalId))
	pkt.WriteByte(p.DestUsrTl)

	for _, d := range p.DestTerminalId {
		pkt.WriteStr(d, numLen) // numLen => cmpp3.0=32字节, cmpp2.0=21字节
	}
	if sp == proto.CMPP3 {
		pkt.WriteByte(p.DestTerminalType)
	}
	pkt.WriteByte(p.MsgLength)
	pkt.WriteStr(p.MsgContent, int(p.MsgLength))

	if sp == proto.CMPP3 {
		pkt.WriteStr(p.LinkId, 20)
	} else {
		pkt.WriteStr(p.LinkId, 8)
	}

	return data
}

// Unpack unpack the binary byte stream to a Cmpp3SubmitReq variable.
// Usually it is used in server side. After unpack, you will get all value of fields in
// Cmpp3SubmitReq struct.
func (p *CmppSubmitReq) Unpack(data []byte, sp proto.SmsProto) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = r.(error)
		}
	}()
	numLen := 32
	if sp == proto.CMPP2 {
		numLen = 21
	}
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

	p.FeeTerminalId = pkt.ReadStr(numLen)
	if sp == proto.CMPP3 {
		p.FeeTerminalType = pkt.ReadByte()
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

	p.DestUsrTl = pkt.ReadByte()

	for i := 0; i < int(p.DestUsrTl); i++ {
		p.DestTerminalId = append(p.DestTerminalId, pkt.ReadStr(numLen))
	}
	if sp == proto.CMPP3 {
		p.DestTerminalType = pkt.ReadByte()
	}
	p.MsgLength = pkt.ReadByte()
	// 0：ASCII 码；3：短信写卡操作；4：二进制信息；8：UCS2 编码；15：含 GBK 汉字。【1字节】
	if p.MsgFmt == 8 {
		p.MsgContent = pkt.ReadUCS2(int(p.MsgLength))
	} else {
		p.MsgContent = pkt.ReadStr(int(p.MsgLength))
	}
	if sp == proto.CMPP3 {
		p.LinkId = pkt.ReadStr(20)
	} else {
		p.LinkId = pkt.ReadStr(8)
	}
	return nil
}
func (p *CmppSubmitReq) Event() event.SmsEvent {
	return event.SmsEventSubmitReq
}
func (p *CmppSubmitReq) SeqId() uint32 {
	return p.seqId
}

// Pack packs the Cmpp3SubmitRsp to bytes stream for Server side.
// Before calling Pack, you should initialize a Cmpp3SubmitRsp variable
// with correct field value.
func (p *CmppSubmitRsp) Pack(seqId uint32, sp proto.SmsProto) []byte {
	var pktLen uint32 = CMPP_HEADER_LEN + 8 + 4
	if sp == proto.CMPP2 {
		pktLen = CMPP_HEADER_LEN + 8 + 1
	}
	data := make([]byte, pktLen)
	pkt := proto.NewPacket(data)

	// Pack header
	pkt.WriteU32(pktLen)
	pkt.WriteU32(CMPP_SUBMIT_RESP.ToInt())

	p.seqId = seqId

	pkt.WriteU32(p.seqId)

	// Pack Body
	pkt.WriteU64(p.MsgId)
	if sp == proto.CMPP3 {
		pkt.WriteU32(p.Result)
	} else {
		pkt.WriteByte(byte(p.Result))
	}
	return data
}

// Unpack unpack the binary byte stream to a Cmpp3SubmitRsp variable.
// Usually it is used in client side. After unpack, you will get all value of fields in
// Cmpp3SubmitRsp struct.
func (p *CmppSubmitRsp) Unpack(data []byte, sp proto.SmsProto) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = r.(error)
		}
	}()
	pkt := proto.NewPacket(data)

	// Sequence Id
	p.seqId = pkt.ReadU32()
	p.MsgId = pkt.ReadU64()
	if sp == proto.CMPP3 {
		p.Result = pkt.ReadU32()
	} else {
		p.Result = uint32(pkt.ReadByte())
	}
	return nil
}
func (p *CmppSubmitRsp) Event() event.SmsEvent {
	return event.SmsEventSubmitRsp
}

func (p *CmppSubmitRsp) SeqId() uint32 {
	return p.seqId
}
