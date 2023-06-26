package smgp

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms/proto"
)

const (
	LoginLen     = 4 + 4 + 4 + 8 + 16 + 1 + 4 + 1 // = 42
	LoginRespLen = 4 + 4 + 4 + 4 + 16 + 1         // = 33
)

type SmgpLoginReq struct {
	seqId               uint32  //  序列 ID
	ClientID            string  //  【8字节】客户端用来登录服务器端的用户账号。
	AuthenticatorClient string  //  【16字节】客户端认证码，用来鉴别客户端的合法性。
	LoginMode           byte    //  【1字节】客户端用来登录服务器端的登录类型。
	Timestamp           uint32  //  【4字节】时间戳
	Version             Version //  【1字节】客户端支持的协议版本号
}

type SmgpLoginRsp struct {
	seqId               uint32  //  序列 ID
	Status              Status  // 状态码，4字节
	AuthenticatorServer string  // 认证串，16字节
	Version             Version // 版本，1字节
}

// func NewLogin(ac *proto.AuthConf, seq uint32) *Login {
// 	lo := &Login{}
// 	lo.PacketLength = LoginLen
// 	lo.RequestId = SMGP_LOGIN
// 	lo.SequenceId = seq
// 	lo.clientID = ac.ClientId
// 	lo.loginMode = 2
// 	var ts string
// 	ts, lo.timestamp = utils.Timestamp2()
// 	authMd5 := md5.Sum(bytes.Join([][]byte{
// 		[]byte(ac.ClientId),
// 		make([]byte, 7),
// 		[]byte(ac.SharedSecret),
// 		[]byte(ts),
// 	}, nil))
// 	lo.authenticatorClient = authMd5[:]
// 	lo.Version = Version(ac.Version)
// 	return lo
// }

func (p *SmgpLoginReq) Pack(seqId uint32) []byte {
	data := make([]byte, LoginLen)
	pkt := proto.NewPacket(data)
	pkt.WriteU32(LoginLen)
	pkt.WriteU32(SMGP_LOGIN.ToInt())
	if seqId > 0 {
		p.seqId = seqId
	}
	pkt.WriteU32(p.seqId)
	pkt.WriteStr(p.ClientID, 8)
	pkt.WriteStr(p.AuthenticatorClient, 16)
	pkt.WriteByte(p.LoginMode)
	pkt.WriteU32(p.Timestamp)
	pkt.WriteByte(byte(p.Version))
	logrus.Warningln("seqid:", seqId, data)
	return data
}

func (p *SmgpLoginReq) Unpack(data []byte) proto.Packer {
	pkt := proto.NewPacket(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	p.ClientID = pkt.ReadStr(8)
	p.AuthenticatorClient = pkt.ReadStr(16)
	p.LoginMode = pkt.ReadByte()
	p.Timestamp = pkt.ReadU32()
	p.Version = Version(pkt.ReadByte())
	return p
}
func (p *SmgpLoginReq) SeqId() uint32 {
	return p.seqId
}

func (p *SmgpLoginReq) String() string {
	return fmt.Sprintf("{SeqId: %s, clientID: %s, authenticatorClient: %x, logoinMode: %x, timestamp: %010d, version: %s}",
		p.seqId, p.ClientID, p.AuthenticatorClient, p.LoginMode, p.Timestamp, p.Version)
}

func (p *SmgpLoginRsp) Pack(seqId uint32) []byte {
	data := make([]byte, LoginLen)
	pkt := proto.NewPacket(data)

	pkt.WriteU32(LoginRespLen)
	pkt.WriteU32(SMGP_LOGIN_RESP.ToInt())
	if seqId > 0 {
		p.seqId = seqId
	}
	pkt.WriteU32(p.seqId)
	pkt.WriteU32(uint32(p.Status))
	pkt.WriteStr(p.AuthenticatorServer, 16)
	pkt.WriteByte(byte(p.Version))
	logrus.Warningln("SmgpLoginRsp.seqid:", seqId, data)
	return data
}

func (p *SmgpLoginRsp) Unpack(data []byte) proto.Packer {
	pkt := proto.NewPacket(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	p.Status = Status(pkt.ReadU32())
	p.AuthenticatorServer = pkt.ReadStr(16)
	p.Version = Version(pkt.ReadByte())
	return p
}
func (p *SmgpLoginRsp) SeqId() uint32 {
	return p.seqId
}

func (p *SmgpLoginRsp) String() string {
	return fmt.Sprintf("{ Header: %s, status: \"%s\", authenticatorISMG: %x, version: %s }",
		p.seqId, p.Status, p.AuthenticatorServer, p.Version)
}
