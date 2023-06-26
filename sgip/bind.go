package sgip

import "github.com/zhiyin2021/zysms/proto"

const SgipBindPkgLen = 61
const SgipBindRspPkgLen = 21

// Bind 登录报文结构体【61 bytes】
type SgipBindReq struct {
	SeqId         []uint32 // 【12 bytes】序列 ID
	LoginType     byte     // 【1 bytes 】登录类型。 1:SP 向 SMG 建立的连接，用于发送命令 2:SMG 向 SP 建立的连接，用于发送命令
	LoginName     string   // 【16 bytes】服务器端给客户端分配的登录名
	LoginPassword string   // 【16 bytes】服务器端和 Login Name 对应的密码
	Reserve       string   // 【8 bytes 】保留字段
}

type SgipBindRsp struct {
	SeqId  []uint32 // 【12 bytes】序列 ID
	Status Status
}

// func NewBind(ac *proto.AuthConf, loginType byte) *Bind {
// 	b := &Bind{}
// 	b.CommandId = SGIP_BIND
// 	b.PacketLength = BindPkgLen
// 	b.SeqId = Sequencer.NextVal()
// 	b.LoginType = loginType
// 	b.LoginName = ac.ClientId
// 	b.LoginPassword = ac.SharedSecret
// 	return b
// }

func (b *SgipBindReq) Pack(seqId []uint32) []byte {
	data := make([]byte, SgipBindPkgLen)
	pkt := proto.NewPacket(data)
	pkt.WriteU32(SgipBindPkgLen)
	pkt.WriteU32(SGIP_BIND.ToInt())
	pkt.WriteU32(seqId[0])
	pkt.WriteU32(seqId[1])
	pkt.WriteU32(seqId[2])
	pkt.WriteByte(b.LoginType)
	pkt.WriteStr(b.LoginName, 16)
	pkt.WriteStr(b.LoginPassword, 16)
	pkt.WriteStr(b.Reserve, 8)
	return data
}

func (b *SgipBindReq) Unpack(data []byte) proto.Packer {
	var pkt = proto.NewPacket(data)
	// Sequence Id
	b.SeqId = make([]uint32, 3)
	b.SeqId[0] = pkt.ReadU32()
	b.SeqId[1] = pkt.ReadU32()
	b.SeqId[2] = pkt.ReadU32()
	b.LoginType = pkt.ReadByte()
	b.LoginName = pkt.ReadStr(16)
	b.LoginPassword = pkt.ReadStr(16)
	b.Reserve = pkt.ReadStr(8)
	return b
}

func (r *SgipBindRsp) Pack(seqId []uint32) []byte {
	data := make([]byte, SgipBindRspPkgLen)
	pkt := proto.NewPacket(data)
	pkt.WriteU32(SgipBindRspPkgLen)
	pkt.WriteU32(SGIP_BIND_RESP.ToInt())
	pkt.WriteU32(seqId[0])
	pkt.WriteU32(seqId[1])
	pkt.WriteU32(seqId[2])
	pkt.WriteByte(byte(r.Status))
	return data
}

func (r *SgipBindRsp) Unpack(data []byte) proto.Packer {
	var pkt = proto.NewPacket(data)
	// Sequence Id
	r.SeqId = make([]uint32, 3)
	r.SeqId[0] = pkt.ReadU32()
	r.SeqId[1] = pkt.ReadU32()
	r.SeqId[2] = pkt.ReadU32()
	r.Status = Status(pkt.ReadByte())
}