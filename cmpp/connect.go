package cmpp

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"

	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/utils"
)

type ConnReq struct {
	base
	SrcAddr string // +6 = 18：源地址，此处为 SP_Id
	AuthSrc string // +16 = 34： 用于鉴别源地址。其值通过单向 MD5 hash 计算得出，表示如下: authenticatorSource = MD5(Source_Addr+9 字节的 0 +shared secret+timestamp) Shared secret 由中国移动与源地址实 体事先商定，timestamp 格式为: MMDDHHMMSS，即月日时分秒，10 位。
	//Version   Version // +1 = 35：双方协商的版本号(高位 4bit 表示主 版本号,低位 4bit 表示次版本号)，对 于3.0的版本，高4bit为3，低4位为 0
	Timestamp uint32 // +4 = 39：时间戳的明文,由客户端产生,格式为 MMDDHHMMSS，即月日时分秒，10 位数字的整型，右对齐。
	Secret    string //非协议内容，调用Pack前需设置
}

type ConnResp struct {
	base
	Status   uint32 // (cmpp3 = 4字节, cmpp2 = 1字节) 0：正确 1：消息结构错 2：非法源地址 3：认证错 4：版本太高 5~ ：其他错误
	AuthIsmg string // 16字节 ISMG认证码，用于鉴别ISMG。其值通过单向MD5 hash计算得出，表示如下： AuthenticatorISMG =MD5（Status+AuthenticatorSource+shared secret），Shared secret 由中国移动与源地址实体事先商定，AuthenticatorSource为源地址实体发送给ISMG的对应消息CMPP_Connect中的值。 认证出错时，此项为空。
	//Version  Version // 1字节 服务器支持的最高版本号，对于3.0的版本，高4bit为3，低4位为0
	Secret  string // 非协议内容
	AuthSrc string // 非协议内容
}

func NewConnReq(ver codec.Version) codec.PDU {
	return &ConnReq{
		base: newBase(ver, CMPP_CONNECT, 0),
	}
}
func NewConnResp(ver codec.Version) codec.PDU {
	return &ConnResp{
		base: newBase(ver, CMPP_CONNECT_RESP, 0),
	}
}

// Pack packs the ConnReq to bytes stream for client side.
// Before calling Pack, you should initialize a ConnReq variable
// with correct SourceAddr(SrcAddr), Secret and Version.
func (p *ConnReq) Marshal(w *codec.BytesWriter) {
	var ts string
	if p.Timestamp == 0 {
		ts, p.Timestamp = utils.Timestamp2() //default: current time.
	} else {
		ts = utils.Timestamp2Str(p.Timestamp)
	}
	md5 := md5.Sum(bytes.Join([][]byte{[]byte(p.SrcAddr),
		make([]byte, 9),
		[]byte(p.Secret),
		[]byte(ts)},
		nil))
	p.AuthSrc = string(md5[:])

	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteStr(p.SrcAddr, 6)
		bw.WriteStr(p.AuthSrc, 16)
		bw.WriteByte(byte(p.Version))
		bw.WriteU32(p.Timestamp)
	})
}

// Unpack unpack the binary byte stream to a ConnReq variable.
// Usually it is used in server side. After unpack, you will get SeqId, SourceAddr,
// AuthenticatorSource, Version and Timestamp.
func (p *ConnReq) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.SrcAddr = br.ReadStr(6)
		p.AuthSrc = br.ReadStr(16)
		p.Version = codec.Version(br.ReadByte())
		p.Timestamp = br.ReadU32()
		return br.Err()
	})
}
func (p *ConnReq) GetResponse() codec.PDU {
	return &ConnResp{
		base: newBase(p.Version, CMPP_CONNECT_RESP, p.SequenceNumber),
	}
}

func (p *ConnResp) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bs := make([]byte, 4)
		binary.BigEndian.PutUint32(bs, p.Status)
		var hash [16]byte
		if p.Version == V30 {
			bw.WriteU32(p.Status)
			hash = md5.Sum(bytes.Join([][]byte{bs,
				[]byte(p.AuthSrc),
				[]byte(p.Secret)},
				nil))
		} else {
			bw.WriteByte(bs[3])
			hash = md5.Sum(bytes.Join([][]byte{
				{bs[3]},
				[]byte(p.AuthSrc),
				[]byte(p.Secret)},
				nil))
		}
		p.AuthIsmg = string(hash[:])
		bw.WriteStr(p.AuthIsmg, 16)
		bw.WriteByte(byte(p.Version))
	})
}

func (p *ConnResp) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		if p.Version == V30 || br.Len() == 21 {
			p.Status = br.ReadU32()
		} else {
			p.Status = uint32(br.ReadByte())
		}
		p.AuthIsmg = br.ReadStr(16)
		p.Version = codec.Version(br.ReadByte())
		return br.Err()
	})
}

func (p *ConnResp) GetResponse() codec.PDU {
	return nil
}
