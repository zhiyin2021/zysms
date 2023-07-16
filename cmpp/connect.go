package cmpp

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"errors"

	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/event"
	"github.com/zhiyin2021/zysms/utils"
)

// Packet length const for cmpp connect request and response packets.
const (
	CmppConnReqLen  uint32 = 4 + 4 + 4 + 6 + 16 + 1 + 4 //39d, 0x27
	Cmpp2ConnRspLen uint32 = 4 + 4 + 4 + 1 + 16 + 1     //30d, 0x1e
	Cmpp3ConnRspLen uint32 = 4 + 4 + 4 + 4 + 16 + 1     //33d, 0x21
)

// Errors for connect resp status.
var (
	ErrnoConnInvalidStruct  uint8 = 1
	ErrnoConnInvalidSrcAddr uint8 = 2
	ErrnoConnAuthFailed     uint8 = 3
	ErrnoConnVerTooHigh     uint8 = 4
	ErrnoConnOthers         uint8 = 5

	ConnRspStatusErrMap = map[uint8]error{
		ErrnoConnInvalidStruct:  errConnInvalidStruct,
		ErrnoConnInvalidSrcAddr: errConnInvalidSrcAddr,
		ErrnoConnAuthFailed:     errConnAuthFailed,
		ErrnoConnVerTooHigh:     errConnVerTooHigh,
		ErrnoConnOthers:         errConnOthers,
	}

	errConnInvalidStruct  = errors.New("connect response status: invalid protocol structure")
	errConnInvalidSrcAddr = errors.New("connect response status: invalid source address")
	errConnAuthFailed     = errors.New("connect response status: auth failed")
	errConnVerTooHigh     = errors.New("connect response status: protocol version is too high")
	errConnOthers         = errors.New("connect response status: other errors")
)

type CmppConnReq struct {
	SrcAddr   string  // +6 = 18：源地址，此处为 SP_Id
	AuthSrc   string  // +16 = 34： 用于鉴别源地址。其值通过单向 MD5 hash 计算得出，表示如下: authenticatorSource = MD5(Source_Addr+9 字节的 0 +shared secret+timestamp) Shared secret 由中国移动与源地址实 体事先商定，timestamp 格式为: MMDDHHMMSS，即月日时分秒，10 位。
	Version   Version // +1 = 35：双方协商的版本号(高位 4bit 表示主 版本号,低位 4bit 表示次版本号)，对 于3.0的版本，高4bit为3，低4位为 0
	Timestamp uint32  // +4 = 39：时间戳的明文,由客户端产生,格式为 MMDDHHMMSS，即月日时分秒，10 位数字的整型，右对齐。
	Secret    string  //非协议内容，调用Pack前需设置
	seqId     uint32  // 序列编号
}

type CmppConnRsp struct {
	Status   uint32  // (cmpp3 = 4字节, cmpp2 = 1字节) 0：正确 1：消息结构错 2：非法源地址 3：认证错 4：版本太高 5~ ：其他错误
	AuthIsmg string  // 16字节 ISMG认证码，用于鉴别ISMG。其值通过单向MD5 hash计算得出，表示如下： AuthenticatorISMG =MD5（Status+AuthenticatorSource+shared secret），Shared secret 由中国移动与源地址实体事先商定，AuthenticatorSource为源地址实体发送给ISMG的对应消息CMPP_Connect中的值。 认证出错时，此项为空。
	Version  Version // 1字节 服务器支持的最高版本号，对于3.0的版本，高4bit为3，低4位为0
	Secret   string  // 非协议内容
	AuthSrc  string  // 非协议内容
	seqId    uint32  // 序列编号
}

// Pack packs the CmppConnReq to bytes stream for client side.
// Before calling Pack, you should initialize a CmppConnReq variable
// with correct SourceAddr(SrcAddr), Secret and Version.
func (p *CmppConnReq) Pack(seqId uint32, sp codec.SmsProto) []byte {

	pkt := codec.NewWriter(CmppConnReqLen, CMPP_CONNECT.ToInt())
	pkt.WriteU32(seqId)
	p.seqId = seqId
	var ts string
	if p.Timestamp == 0 {
		ts, p.Timestamp = utils.Timestamp2() //default: current time.
	} else {
		ts = utils.Timestamp2Str(p.Timestamp)
	}

	// Pack body
	pkt.WriteStr(p.SrcAddr, 6)

	md5 := md5.Sum(bytes.Join([][]byte{[]byte(p.SrcAddr),
		make([]byte, 9),
		[]byte(p.Secret),
		[]byte(ts)},
		nil))
	p.AuthSrc = string(md5[:])

	pkt.WriteStr(p.AuthSrc, 16)
	pkt.WriteByte(byte(p.Version))
	pkt.WriteU32(p.Timestamp)

	return pkt.Bytes()
}

// Unpack unpack the binary byte stream to a CmppConnReq variable.
// Usually it is used in server side. After unpack, you will get SeqId, SourceAddr,
// AuthenticatorSource, Version and Timestamp.
func (p *CmppConnReq) Unpack(data []byte, sp codec.SmsProto) error {
	pkt := codec.NewReader(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	// Body: Source_Addr
	p.SrcAddr = pkt.ReadStr(6)
	// Body: AuthSrc
	p.AuthSrc = pkt.ReadStr(16)
	// Body: Version
	p.Version = Version(pkt.ReadByte())
	// Body: timestamp
	p.Timestamp = pkt.ReadU32()
	return pkt.Err()
}
func (p *CmppConnReq) Event() event.SmsEvent {
	return event.SmsEventAuthReq
}

func (p *CmppConnReq) SeqId() uint32 {
	return p.seqId
}

func (p *CmppConnRsp) Pack(seqId uint32, sp codec.SmsProto) []byte {
	rspLen := Cmpp2ConnRspLen
	if sp == codec.CMPP30 {
		rspLen = Cmpp3ConnRspLen
	}
	pkt := codec.NewWriter(rspLen, CMPP_CONNECT_RESP.ToInt())
	pkt.WriteU32(seqId)

	p.seqId = seqId

	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, p.Status)

	if sp == codec.CMPP30 {
		// pack body
		pkt.WriteU32(p.Status)
	} else {
		pkt.WriteByte(bs[3])
	}

	hash := md5.Sum(bytes.Join([][]byte{bs,
		[]byte(p.AuthSrc),
		[]byte(p.Secret)},
		nil))
	p.AuthIsmg = string(hash[:])
	pkt.WriteStr(p.AuthIsmg, 16)

	pkt.WriteByte(byte(p.Version))

	return pkt.Bytes()
}

func (p *CmppConnRsp) Unpack(data []byte, sp codec.SmsProto) error {
	pkt := codec.NewReader(data)

	// Sequence Id
	p.seqId = pkt.ReadU32()
	if sp == codec.CMPP30 {
		// Body: Status
		p.Status = pkt.ReadU32()
	} else {
		p.Status = uint32(pkt.ReadByte())
	}

	// Body: AuthenticatorISMG
	p.AuthIsmg = pkt.ReadStr(16)
	// Body: Version
	p.Version = Version(pkt.ReadByte())
	return pkt.Err()
}
func (p *CmppConnRsp) Event() event.SmsEvent {
	return event.SmsEventAuthRsp
}

func (p *CmppConnRsp) SeqId() uint32 {
	return p.seqId
}
