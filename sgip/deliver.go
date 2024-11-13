package sgip

import "github.com/zhiyin2021/zysms/codec"

type DeliverReq struct {
	base
	UserNumber    string             //  接收该短消息的手机号，该字段重复UserCount指定的次数，手机号码前加“86”国别标志【 21 bytes 】
	SPNumber      string             //  SP的接入号码【 21 bytes 】
	TpPid         byte               //  GSM协议类型。详细解释请参考GSM03.40中的9.2.3.9 【 1  bytes 】
	TpUdhi        byte               //  GSM协议类型。详细解释请参考GSM03.40中的9.2.3.9 【 1  bytes 】
	MessageCoding byte               //  短消息的编码格式。 【 1  bytes 】
	Message       codec.ShortMessage //[]byte // 信息内容
	Reserve       string             //  保留，扩展用【 8 bytes 】
}
type DeliverResp struct {
	base
	Status  Status
	Reserve string
}

func NewDeliverReq(ver codec.Version, nodeId uint32) codec.PDU {
	return &DeliverReq{
		base: newBase(ver, SGIP_DELIVER, [3]uint32{nodeId, 0, 0}),
	}
}
func NewDeliverResp(ver codec.Version, nodeId uint32) codec.PDU {
	return &DeliverResp{
		base: newBase(ver, SGIP_DELIVER_RESP, [3]uint32{nodeId, 0, 0}),
	}
}
func (p *DeliverReq) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteStr(p.UserNumber, 21)
		bw.WriteStr(p.SPNumber, 21)
		bw.WriteByte(p.TpPid)
		bw.WriteByte(p.TpUdhi)
		bw.WriteByte(p.MessageCoding)
		p.Message.Marshal(bw)
		bw.WriteStr(p.Reserve, 8)
	})
}

func (p *DeliverReq) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.UserNumber = br.ReadStr(21)
		p.SPNumber = br.ReadStr(21)
		p.TpPid = br.ReadU8()
		p.TpUdhi = br.ReadU8()
		p.MessageCoding = br.ReadU8()
		p.Message.Unmarshal(br, p.TpUdhi == 1, p.MessageCoding)
		p.Reserve = br.ReadStr(8)
		return br.Err()
	})
}

func (b *DeliverReq) GetResponse() codec.PDU {
	return &DeliverResp{
		base: newBase(b.Version, SGIP_DELIVER_RESP, b.SequenceNumber),
	}
}

func (p *DeliverResp) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteByte(byte(p.Status))
		bw.WriteStr(p.Reserve, 8)
	})
}
func (p *DeliverResp) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.Status = Status(br.ReadU8())
		p.Reserve = br.ReadStr(8)
		return br.Err()
	})
}
func (p *DeliverResp) GetResponse() codec.PDU {
	return nil
}
