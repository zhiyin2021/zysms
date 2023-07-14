package cmpp

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms/event"
	"github.com/zhiyin2021/zysms/proto"
	"github.com/zhiyin2021/zysms/utils"
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
	MsgLength          uint8
	MsgContent         []byte // 信息内容
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
	var pktLen uint32 = CMPP_HEADER_LEN + 117 + uint32(len(p.DestTerminalId))*21 + 1 + uint32(p.MsgLength) + 8
	numLen := 21
	if sp == proto.CMPP30 {
		pktLen = CMPP_HEADER_LEN + 129 + uint32(len(p.DestTerminalId)*32) + 1 + 1 + uint32(p.MsgLength) + 20
		numLen = 32
	}
	logrus.Infof("submit.pack %d,%d", pktLen, numLen)

	pkt := proto.NewCmppBuffer(pktLen, CMPP_SUBMIT.ToInt(), seqId)
	// Pack header

	p.seqId = seqId

	// Pack Body
	pkt.WriteU64(p.MsgId)

	if p.PkTotal == 0 && p.PkNumber == 0 {
		p.PkTotal, p.PkNumber = 1, 1
	}
	pkt.WriteByte(p.PkTotal)
	pkt.WriteByte(p.PkNumber)
	pkt.WriteByte(p.RegisteredDelivery)
	pkt.WriteByte(p.MsgLevel)
	pkt.WriteCStrN(p.ServiceId, 10)
	pkt.WriteByte(p.FeeUserType)
	pkt.WriteCStrN(p.FeeTerminalId, numLen) // numLen => cmpp3.0=32字节, cmpp2.0=21字节
	if sp == proto.CMPP30 {
		pkt.WriteByte(p.FeeTerminalType)
	}
	pkt.WriteByte(p.TpPid)
	pkt.WriteByte(p.TpUdhi)
	pkt.WriteByte(p.MsgFmt)
	pkt.WriteCStrN(p.MsgSrc, 6)
	pkt.WriteCStrN(p.FeeType, 2)
	pkt.WriteCStrN(p.FeeCode, 6)
	pkt.WriteCStrN(p.ValidTime, 17)
	pkt.WriteCStrN(p.AtTime, 17)
	pkt.WriteCStrN(p.SrcId, 21)
	p.DestUsrTl = uint8(len(p.DestTerminalId))
	pkt.WriteByte(p.DestUsrTl)

	for _, d := range p.DestTerminalId {
		pkt.WriteCStrN(d, numLen) // numLen => cmpp3.0=32字节, cmpp2.0=21字节
	}
	if sp == proto.CMPP30 {
		pkt.WriteByte(p.DestTerminalType)
	}
	pkt.WriteByte(p.MsgLength)

	pkt.WriteBytes(p.MsgContent)

	if sp == proto.CMPP30 {
		pkt.WriteCStrN(p.LinkId, 20)
	} else {
		pkt.WriteCStrN(p.LinkId, 8)
	}

	return pkt.Bytes()
}
func (p *CmppSubmitReq) ContentText() string {
	// 0：ASCII 码；3：短信写卡操作；4：二进制信息；8：UCS2 编码；15：含 GBK 汉字。【1字节】
	if p.MsgFmt == 8 {
		txt, _ := utils.Ucs2ToUtf8(p.MsgContent)
		return string(txt)
	}
	return string(p.MsgContent)
}

// Unpack unpack the binary byte stream to a Cmpp3SubmitReq variable.
// Usually it is used in server side. After unpack, you will get all value of fields in
// Cmpp3SubmitReq struct.
func (p *CmppSubmitReq) Unpack(data []byte, sp proto.SmsProto) (e error) {
	numLen := 21
	if sp == proto.CMPP30 {
		numLen = 32
	}
	pkt := proto.NewBuffer(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	p.MsgId = pkt.ReadU64()
	p.PkTotal = pkt.ReadByte()
	p.PkNumber = pkt.ReadByte()
	p.RegisteredDelivery = pkt.ReadByte()
	p.MsgLevel = pkt.ReadByte()
	p.ServiceId = pkt.ReadCStrN(10)
	p.FeeUserType = pkt.ReadByte()
	p.FeeTerminalId = pkt.ReadCStrN(numLen)
	if sp == proto.CMPP30 {
		p.FeeTerminalType = pkt.ReadByte()
	}
	p.TpPid = pkt.ReadByte()
	p.TpUdhi = pkt.ReadByte()
	p.MsgFmt = pkt.ReadByte()
	p.MsgSrc = pkt.ReadCStrN(6)
	p.FeeType = pkt.ReadCStrN(2)
	p.FeeCode = pkt.ReadCStrN(6)
	p.ValidTime = pkt.ReadCStrN(17)
	p.AtTime = pkt.ReadCStrN(17)
	p.SrcId = pkt.ReadCStrN(21)
	p.DestUsrTl = pkt.ReadByte()
	for i := 0; i < int(p.DestUsrTl); i++ {
		p.DestTerminalId = append(p.DestTerminalId, pkt.ReadCStrN(numLen))
	}
	if sp == proto.CMPP30 {
		p.DestTerminalType = pkt.ReadByte()
	}
	p.MsgLength = pkt.ReadByte()

	// if p.MsgFmt == 8 {
	// 	p.MsgContent = pkt.ReadUCS2(int(p.MsgLength))
	// } else {
	// 	p.MsgContent = pkt.ReadStr(int(p.MsgLength))
	// }

	p.MsgContent = pkt.ReadN(int(p.MsgLength))
	if sp == proto.CMPP30 {
		p.LinkId = pkt.ReadCStrN(20)
	} else {
		p.LinkId = pkt.ReadCStrN(8)
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
	var pktLen uint32 = CMPP_HEADER_LEN + 8 + 1
	if sp == proto.CMPP30 {
		pktLen = CMPP_HEADER_LEN + 8 + 4
	}
	pkt := proto.NewCmppBuffer(pktLen, CMPP_SUBMIT_RESP.ToInt(), seqId)
	p.seqId = seqId

	// Pack Body
	pkt.WriteU64(p.MsgId)
	if sp == proto.CMPP30 {
		pkt.WriteU32(p.Result)
	} else {
		pkt.WriteByte(byte(p.Result))
	}
	return pkt.Bytes()
}

// Unpack unpack the binary byte stream to a Cmpp3SubmitRsp variable.
// Usually it is used in client side. After unpack, you will get all value of fields in
// Cmpp3SubmitRsp struct.
func (p *CmppSubmitRsp) Unpack(data []byte, sp proto.SmsProto) error {
	pkt := proto.NewBuffer(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	p.MsgId = pkt.ReadU64()
	if sp == proto.CMPP30 {
		p.Result = pkt.ReadU32()
	} else {
		p.Result = uint32(pkt.ReadByte())
	}
	return pkt.Err()
}
func (p *CmppSubmitRsp) Event() event.SmsEvent {
	return event.SmsEventSubmitRsp
}

func (p *CmppSubmitRsp) SeqId() uint32 {
	return p.seqId
}
