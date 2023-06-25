package codec

import (
	"errors"

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

type Cmpp2DeliverReq struct {
	MsgId            uint64 // 信息标识，由 SP 接入的短信网关本身产生，本处填空(0)。【8字节】
	DestId           string // 目的号码 21
	ServiceId        string // 业务类型 10
	TpPid            uint8  // GSM协议类型。详细是解释请参考GSM03.40中的9.2.3.9 【1字节】
	TpUdhi           uint8  // GSM协议类型。详细是解释请参考GSM03.40中的9.2.3.9 【1字节】
	MsgFmt           uint8  // 信息格式 【1字节】
	SrcTerminalId    string // 源终端MSISDN号码 【21字节】
	RegisterDelivery uint8  // 是否要求返回状态确认报告
	MsgLength        uint8  // 信息长度
	MsgContent       string // 信息内容
	Reserve          string // 保留

	//session info
	SeqId uint32 // sequence id
}

type Cmpp2DeliverRsp struct {
	MsgId  uint64
	Result uint8

	//session info
	SeqId uint32
}
type Cmpp3DeliverReq struct {
	MsgId            uint64 // 消息标识
	DestId           string // 目的号码 21
	ServiceId        string // 业务类型 10
	TpPid            uint8  // GSM协议类型。详细是解释请参考GSM03.40中的9.2.3.9 【1字节】
	TpUdhi           uint8  // GSM协议类型。详细是解释请参考GSM03.40中的9.2.3.9 【1字节】
	MsgFmt           uint8  // 信息格式 【1字节】
	SrcTerminalId    string // 源终端MSISDN号码（状态报告时填为CMPP_SUBMIT消息的目的终端号码）
	SrcTerminalType  uint8  // 源终端号码类型，0：真实号码；1：伪码
	RegisterDelivery uint8  // 是否为状态报告
	MsgLength        uint8
	MsgContent       string
	LinkId           string

	//session info
	SeqId uint32
}
type Cmpp3DeliverRsp struct {
	MsgId  uint64
	Result uint32

	//session info
	SeqId uint32
}

// Pack packs the Cmpp2DeliverReq to bytes stream for client side.
func (p *Cmpp2DeliverReq) Pack(seqId uint32) []byte {
	pktLen := CMPP_HEADER_LEN + 65 + uint32(p.MsgLength) + 8
	data := make([]byte, pktLen)
	pkt := proto.NewPacket(data)
	// Pack header
	pkt.WriteU32(pktLen)
	pkt.WriteU32(CMPP_DELIVER.ToInt())
	pkt.WriteU32(seqId)
	p.SeqId = seqId

	// Pack Body
	pkt.WriteU64(p.MsgId)
	pkt.WriteStr(p.DestId, 21)
	pkt.WriteStr(p.ServiceId, 10)
	pkt.WriteByte(p.TpPid)
	pkt.WriteByte(p.TpUdhi)
	pkt.WriteByte(p.MsgFmt)
	pkt.WriteStr(p.SrcTerminalId, 21)
	pkt.WriteByte(p.RegisterDelivery)
	pkt.WriteByte(p.MsgLength)
	pkt.WriteStr(p.MsgContent, int(p.MsgLength))
	pkt.WriteStr(p.Reserve, 8)

	return data
}

// Unpack unpack the binary byte stream to a Cmpp2DeliverReq variable.
// After unpack, you will get all value of fields in
// Cmpp2DeliverReq struct.
func (p *Cmpp2DeliverReq) Unpack(data []byte) {

	pkt := proto.NewPacket(data)

	// Sequence Id
	p.SeqId = pkt.ReadU32()

	// Body
	p.MsgId = pkt.ReadU64()

	p.DestId = pkt.ReadStr(21)

	p.ServiceId = pkt.ReadStr(10)

	p.TpPid = pkt.ReadByte()
	p.TpUdhi = pkt.ReadByte()
	p.MsgFmt = pkt.ReadByte()

	p.SrcTerminalId = pkt.ReadStr(21)

	p.RegisterDelivery = pkt.ReadByte()
	p.MsgLength = pkt.ReadByte()

	p.MsgContent = pkt.ReadStr(int(p.MsgLength))

	p.Reserve = pkt.ReadStr(8)

}

// Pack packs the Cmpp2DeliverRsp to bytes stream for client side.
func (p *Cmpp2DeliverRsp) Pack(seqId uint32) []byte {
	data := make([]byte, Cmpp2DeliverRspLen)
	pkt := proto.NewPacket(data)

	// Pack header
	pkt.WriteU32(Cmpp2DeliverRspLen)
	pkt.WriteU32(CMPP_DELIVER_RESP.ToInt())
	pkt.WriteU32(seqId)
	p.SeqId = seqId

	// Pack Body
	pkt.WriteU64(p.MsgId)
	pkt.WriteByte(p.Result)
	return data
}

// Unpack unpack the binary byte stream to a Cmpp2DeliverRsp variable.
// After unpack, you will get all value of fields in
// Cmpp2DeliverRsp struct.
func (p *Cmpp2DeliverRsp) Unpack(data []byte) {
	pkt := proto.NewPacket(data)
	// Sequence Id
	p.SeqId = pkt.ReadU32()
	p.MsgId = pkt.ReadU64()
	p.Result = pkt.ReadByte()
}

// Pack packs the Cmpp3DeliverReq to bytes stream for client side.
func (p *Cmpp3DeliverReq) Pack(seqId uint32) []byte {
	var pktLen uint32 = CMPP_HEADER_LEN + 77 + uint32(p.MsgLength) + 20
	data := make([]byte, pktLen)
	pkt := proto.NewPacket(data)

	// Pack header
	pkt.WriteU32(pktLen)
	pkt.WriteU32(CMPP_DELIVER.ToInt())
	pkt.WriteU32(seqId)
	p.SeqId = seqId

	// Pack Body
	pkt.WriteU64(p.MsgId)
	pkt.WriteStr(p.DestId, 21)
	pkt.WriteStr(p.ServiceId, 10)
	pkt.WriteByte(p.TpPid)
	pkt.WriteByte(p.TpUdhi)
	pkt.WriteByte(p.MsgFmt)
	pkt.WriteStr(p.SrcTerminalId, 32)
	pkt.WriteByte(p.SrcTerminalType)
	pkt.WriteByte(p.RegisterDelivery)
	pkt.WriteByte(p.MsgLength)
	pkt.WriteStr(p.MsgContent, int(p.MsgLength))
	pkt.WriteStr(p.LinkId, 20)

	return data
}

// Unpack unpack the binary byte stream to a Cmpp3DeliverReq variable.
// After unpack, you will get all value of fields in
// Cmpp3DeliverReq struct.
func (p *Cmpp3DeliverReq) Unpack(data []byte) {
	pkt := proto.NewPacket(data)

	// Sequence Id
	p.SeqId = pkt.ReadU32()

	// Body
	p.MsgId = pkt.ReadU64()

	p.DestId = pkt.ReadStr(21)

	p.ServiceId = pkt.ReadStr(10)
	p.TpPid = pkt.ReadByte()
	p.TpUdhi = pkt.ReadByte()
	p.MsgFmt = pkt.ReadByte()

	p.SrcTerminalId = pkt.ReadStr(32)
	p.SrcTerminalType = pkt.ReadByte()

	p.RegisterDelivery = pkt.ReadByte()
	p.MsgLength = pkt.ReadByte()

	p.MsgContent = pkt.ReadStr(int(p.MsgLength))

	p.LinkId = pkt.ReadStr(20)
}

// Pack packs the Cmpp3DeliverRsp to bytes stream for client side.
func (p *Cmpp3DeliverRsp) Pack(seqId uint32) []byte {
	data := make([]byte, Cmpp3DeliverRspLen)
	pkt := proto.NewPacket(data)

	// Pack header
	pkt.WriteU32(Cmpp3DeliverRspLen)
	pkt.WriteU32(CMPP_DELIVER_RESP.ToInt())
	pkt.WriteU32(seqId)
	p.SeqId = seqId

	// Pack Body
	pkt.WriteU64(p.MsgId)
	pkt.WriteU32(p.Result)

	return data
}

// Unpack unpack the binary byte stream to a Cmpp3DeliverRsp variable.
// After unpack, you will get all value of fields in
// Cmpp3DeliverRsp struct.
func (p *Cmpp3DeliverRsp) Unpack(data []byte) {
	pkt := proto.NewPacket(data)
	// Sequence Id
	p.SeqId = pkt.ReadU32()
	p.MsgId = pkt.ReadU64()
	p.Result = pkt.ReadU32()
}
