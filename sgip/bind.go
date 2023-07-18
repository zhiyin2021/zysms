package sgip

import "github.com/zhiyin2021/zysms/codec"

// Bind 登录报文结构体【61 bytes】
type BindReq struct {
	base
	LoginType     byte   // 【1 bytes 】登录类型。 1:SP 向 SMG 建立的连接，用于发送命令 2:SMG 向 SP 建立的连接，用于发送命令
	LoginName     string // 【16 bytes】服务器端给客户端分配的登录名
	LoginPassword string // 【16 bytes】服务器端和 Login Name 对应的密码
	Reserve       string // 【8 bytes 】保留字段
}

type BindResp struct {
	base
	Status Status
}

func NewBindReq(ver codec.Version, nodeId uint32) codec.PDU {
	return &BindReq{
		base: newBase(ver, SGIP_BIND, [3]uint32{nodeId, 0, 0}),
	}
}
func NewBindResp(ver codec.Version, nodeId uint32) codec.PDU {
	return &BindResp{
		base: newBase(ver, SGIP_BIND_RESP, [3]uint32{nodeId, 0, 0}),
	}
}

func (p *BindReq) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteByte(p.LoginType)
		bw.WriteStr(p.LoginName, 16)
		bw.WriteStr(p.LoginPassword, 16)
		bw.WriteStr(p.Reserve, 8)
	})
}

func (p *BindReq) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.LoginType = br.ReadByte()
		p.LoginName = br.ReadStr(16)
		p.LoginPassword = br.ReadStr(16)
		p.Reserve = br.ReadStr(8)
		return br.Err()
	})
}

func (b *BindReq) GetResponse() codec.PDU {
	return &BindResp{
		base: newBase(b.Version, SGIP_BIND_RESP, b.SequenceNumber),
	}
}

func (p *BindResp) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteByte(byte(p.Status))
	})
}

func (p *BindResp) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.Status = Status(br.ReadByte())
		return br.Err()
	})
}
func (p *BindResp) GetResponse() codec.PDU {
	return nil
}
