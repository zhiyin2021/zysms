package smgp

import (
	"bytes"
	"crypto/md5"

	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/utils"
)

type LoginReq struct {
	base
	ClientID            string //  【8字节】客户端用来登录服务器端的用户账号。
	AuthenticatorClient string //  【16字节】客户端认证码，用来鉴别客户端的合法性。
	LoginMode           byte   //  【1字节】客户端用来登录服务器端的登录类型。
	Timestamp           uint32 //  【4字节】时间戳
	Secret              string //非协议内容，调用Pack前需设置
	// Version             codec.Version //  【1字节】客户端支持的协议版本号
}

type LoginResp struct {
	base
	Status              Status // 状态码，4字节
	AuthenticatorServer string // 认证串，16字节
	// Version             Version // 版本，1字节
}

func NewLoginReq(ver codec.Version) codec.PDU {
	return &LoginReq{
		base:      newBase(ver, SMGP_LOGIN, 0),
		LoginMode: 2,
	}
}

func NewLoginResp(ver codec.Version) codec.PDU {
	return &LoginResp{
		base: newBase(ver, SMGP_LOGIN_RESP, 0),
	}
}

// Pack packs the ActiveTestReq to bytes stream for client side.
func (p *LoginReq) Marshal(w *codec.BytesWriter) {
	var ts string
	if p.Timestamp == 0 {
		ts, p.Timestamp = utils.Timestamp2() //default: current time.
	} else {
		ts = utils.Timestamp2Str(p.Timestamp)
	}
	md5 := md5.Sum(bytes.Join([][]byte{[]byte(p.ClientID),
		make([]byte, 7),
		[]byte(p.Secret),
		[]byte(ts)},
		nil))
	p.AuthenticatorClient = string(md5[:])

	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteStr(p.ClientID, 8)
		bw.WriteStr(p.AuthenticatorClient, 16)
		bw.WriteByte(p.LoginMode)
		bw.WriteU32(p.Timestamp)
		bw.WriteByte(byte(p.Version))
	})
}

// Unpack unpack the binary byte stream to a ActiveTestReq variable.
// After unpack, you will get all value of fields in
func (p *LoginReq) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.ClientID = br.ReadStr(8)
		p.AuthenticatorClient = br.ReadStr(16)
		p.LoginMode = br.ReadU8()
		p.Timestamp = br.ReadU32()
		p.Version = codec.Version(br.ReadU8())
		return br.Err()
	})
}

// GetResponse implements PDU interface.
func (b *LoginReq) GetResponse() codec.PDU {
	return &LoginResp{
		base: newBase(b.Version, SMGP_LOGIN_RESP, b.SequenceNumber),
	}
}

// Pack packs the ActiveTestReq to bytes stream for client side.
func (p *LoginResp) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteU32(uint32(p.Status))
		bw.WriteStr(p.AuthenticatorServer, 16)
		bw.WriteByte(byte(p.Version))
	})
}

// Unpack unpack the binary byte stream to a ActiveTestReq variable.
// After unpack, you will get all value of fields in
func (p *LoginResp) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.Status = Status(br.ReadU32())
		p.AuthenticatorServer = br.ReadStr(16)
		p.Version = codec.Version(br.ReadU8())
		return br.Err()
	})
}
func (p *LoginResp) GetResponse() codec.PDU {
	return nil
}
