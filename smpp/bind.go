package smpp

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"

	"github.com/zhiyin2021/zysms/event"
	"github.com/zhiyin2021/zysms/proto"
	"github.com/zhiyin2021/zysms/utils"
)

// Packet length const for cmpp connect request and response packets.
const (
	CmppConnReqLen  uint32 = 4 + 4 + 4 + 6 + 16 + 1 + 4 //39d, 0x27
	Cmpp2ConnRspLen uint32 = 4 + 4 + 4 + 1 + 16 + 1     //30d, 0x1e
	Cmpp3ConnRspLen uint32 = 4 + 4 + 4 + 4 + 16 + 1     //33d, 0x21
)

type SmppBindReq struct {
	// header 消息头 12字节
	SrcAddr   string  // +6 = 18：源地址，此处为 SP_Id
	AuthSrc   string  // +16 = 34： 用于鉴别源地址。其值通过单向 MD5 hash 计算得出，表示如下: authenticatorSource = MD5(Source_Addr+9 字节的 0 +shared secret+timestamp) Shared secret 由中国移动与源地址实 体事先商定，timestamp 格式为: MMDDHHMMSS，即月日时分秒，10 位。
	Version   Version // +1 = 35：双方协商的版本号(高位 4bit 表示主 版本号,低位 4bit 表示次版本号)，对 于3.0的版本，高4bit为3，低4位为 0
	Timestamp uint32  // +4 = 39：时间戳的明文,由客户端产生,格式为 MMDDHHMMSS，即月日时分秒，10 位数字的整型，右对齐。
	Secret    string  //非协议内容，调用Pack前需设置
	seqId     uint32  // 序列编号
}

type SmppBindRsp struct {
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
func (p *CmppConnReq) Pack(seqId uint32, sp proto.SmsProto) []byte {
	buf := make([]byte, CmppConnReqLen)
	pkt := proto.NewPacket(buf)

	// Pack header
	pkt.WriteU32(CmppConnReqLen)
	pkt.WriteU32(CMPP_CONNECT.ToInt())

	p.seqId = seqId

	pkt.WriteU32(p.seqId)

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

	return buf
}

// Unpack unpack the binary byte stream to a CmppConnReq variable.
// Usually it is used in server side. After unpack, you will get SeqId, SourceAddr,
// AuthenticatorSource, Version and Timestamp.
func (p *CmppConnReq) Unpack(data []byte, sp proto.SmsProto) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = r.(error)
		}
	}()
	pkt := proto.NewPacket(data)
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
	return nil
}
func (p *CmppConnReq) Event() event.SmsEvent {
	return event.SmsEventAuthReq
}

func (p *CmppConnReq) SeqId() uint32 {
	return p.seqId
}

func (p *CmppConnRsp) Pack(seqId uint32, sp proto.SmsProto) []byte {
	rspLen := Cmpp2ConnRspLen
	if sp == proto.CMPP30 {
		rspLen = Cmpp3ConnRspLen
	}
	data := make([]byte, rspLen)
	pkt := proto.NewPacket(data)

	// pack header
	pkt.WriteU32(rspLen)

	pkt.WriteU32(CMPP_CONNECT_RESP.ToInt())
	p.seqId = seqId

	pkt.WriteU32(p.seqId)

	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, p.Status)

	if sp == proto.CMPP30 {
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

	return data
}

func (p *CmppConnRsp) Unpack(data []byte, sp proto.SmsProto) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = r.(error)
		}
	}()
	pkt := proto.NewPacket(data)

	// Sequence Id
	p.seqId = pkt.ReadU32()
	if sp == proto.CMPP30 {
		// Body: Status
		p.Status = pkt.ReadU32()
	} else {
		p.Status = uint32(pkt.ReadByte())
	}

	// Body: AuthenticatorISMG
	p.AuthIsmg = pkt.ReadStr(16)
	// Body: Version
	p.Version = Version(pkt.ReadByte())
	return nil
}
func (p *CmppConnRsp) Event() event.SmsEvent {
	return event.SmsEventAuthRsp
}

func (p *CmppConnRsp) SeqId() uint32 {
	return p.seqId
}
