package sgip

import (
	"github.com/zhiyin2021/zysms/proto"
)

type SgipDeliverReq struct {
	SeqId          []uint32 // 消息流水号，由SP接入的短信网关本身产生  12 bytes 】
	UserNumber     string   //  接收该短消息的手机号，该字段重复UserCount指定的次数，手机号码前加“86”国别标志【 21 bytes 】
	SPNumber       string   //  SP的接入号码【 21 bytes 】
	TpPid          byte     //  GSM协议类型。详细解释请参考GSM03.40中的9.2.3.9 【 1  bytes 】
	TpUdhi         byte     //  GSM协议类型。详细解释请参考GSM03.40中的9.2.3.9 【 1  bytes 】
	MessageCoding  byte     //  短消息的编码格式。 【 1  bytes 】
	MessageLength  uint32   //  短消息的长度【 4 bytes 】
	MessageContent string   //  编码后消息内容
	Reserve        string   //  保留，扩展用【 8 bytes 】
}
type SgipDeliverRsp struct {
	SeqId   []uint32 // 消息流水号，由SP接入的短信网关本身产生  12 bytes 】
	Status  Status
	Reserve string
}

const (
	SgipDeliverReqLen = SGIP_HEADER_LEN + 57 // SGIP_HEADER_LEN + 57 + MessageLength
	SgipDeliverRspLen = SGIP_HEADER_LEN + 9
)

// func NewDeliver(ac *proto.AuthConf, phone, content, destNo string) proto.RequestPdu {
// 	dlv := &Deliver{}
// 	dlv.PacketLength = MoBaseLen
// 	dlv.CommandId = SGIP_DELIVER
// 	dlv.SequenceNumber = Sequencer.NextVal()
// 	dlv.UserNumber = phone
// 	if !strings.HasPrefix(destNo, ac.SmsDisplayNo) {
// 		destNo = ac.SmsDisplayNo + destNo
// 	}
// 	dlv.SPNumber = destNo
// 	dlv.MessageCoding = utils.MsgFmt(content)

// 	var bs []byte
// 	// 上行短信不支持长短信，过长内容会被截取
// 	if dlv.MessageCoding == 8 {
// 		bs, _ = utils.Utf8ToUcs2(content)
// 		if len(bs) > 140 {
// 			bs = bs[:140]
// 		}
// 	} else {
// 		bs = []byte(content)
// 		if len(bs) > 160 {
// 			bs = bs[:160]
// 		}
// 	}
// 	dlv.MessageLength = uint32(len(bs))
// 	dlv.MessageContent = bs
// 	dlv.PacketLength = MoBaseLen + dlv.MessageLength
// 	return dlv
// }

func (d *SgipDeliverReq) Pack(seqId []uint32) []byte {
	pktLen := SgipDeliverReqLen + uint32(d.MessageLength)
	data := make([]byte, pktLen)
	pkt := proto.NewPacket(data)
	pkt.WriteU32(pktLen)
	pkt.WriteU32(SGIP_DELIVER.ToInt())
	pkt.WriteU32(seqId[0])
	pkt.WriteU32(seqId[1])
	pkt.WriteU32(seqId[2])
	pkt.WriteStr(d.UserNumber, 21)
	pkt.WriteStr(d.SPNumber, 21)
	pkt.WriteByte(d.TpPid)
	pkt.WriteByte(d.TpUdhi)
	pkt.WriteByte(d.MessageCoding)
	pkt.WriteU32(d.MessageLength)
	pkt.WriteStr(d.MessageContent, int(d.MessageLength))
	pkt.WriteStr(d.Reserve, 8)
	return data
}

func (d *SgipDeliverReq) Unpack(data []byte) proto.Packer {
	var pkt = proto.NewPacket(data)
	// Sequence Id
	d.SeqId = make([]uint32, 3)
	d.SeqId[0] = pkt.ReadU32()
	d.SeqId[1] = pkt.ReadU32()
	d.SeqId[2] = pkt.ReadU32()
	d.UserNumber = pkt.ReadStr(21)
	d.SPNumber = pkt.ReadStr(21)
	d.TpPid = pkt.ReadByte()
	d.TpUdhi = pkt.ReadByte()
	d.MessageCoding = pkt.ReadByte()
	d.MessageLength = pkt.ReadU32()
	d.MessageContent = pkt.ReadStr(int(d.MessageLength))
	d.Reserve = pkt.ReadStr(8)
}

// func (d *Deliver) Log() []log.Field {
// 	ls := d.MessageHeader.Log()
// 	var l = len(d.MessageContent)
// 	if l > 6 {
// 		l = 6
// 	}
// 	msg := hex.EncodeToString(d.MessageContent[:l]) + "..."
// 	return append(ls,
// 		log.String("userNumber", d.UserNumber),
// 		log.String("spNumber", d.SPNumber),
// 		log.Uint8("msgFormat", d.MessageCoding),
// 		log.Uint32("msgLength", d.MessageLength),
// 		log.String("msgContent", msg),
// 		log.Uint8("tpPid", d.TpPid),
// 		log.Uint8("tpUdhi", d.TpUdhi),
// 	)
// }

func (r *SgipDeliverRsp) Pack(seqId []uint32) []byte {
	data := make([]byte, SgipDeliverRspLen)
	pkt := proto.NewPacket(data)
	pkt.WriteU32(SgipDeliverRspLen)
	pkt.WriteU32(SGIP_DELIVER_RESP.ToInt())
	pkt.WriteU32(seqId[0])
	pkt.WriteU32(seqId[1])
	pkt.WriteU32(seqId[2])
	pkt.WriteU32(uint32(r.Status))
	pkt.WriteStr(r.Reserve, 8)
	return data
}

func (r *SgipDeliverRsp) Unpack(data []byte) proto.Packer {
	var pkt = proto.NewPacket(data)
	// Sequence Id
	r.SeqId = make([]uint32, 3)
	r.SeqId[0] = pkt.ReadU32()
	r.SeqId[1] = pkt.ReadU32()
	r.SeqId[2] = pkt.ReadU32()
	r.Status = Status(pkt.ReadU32())
	r.Reserve = pkt.ReadStr(8)
}

// func (r *DeliverRsp) Log() []log.Field {
// 	ls := r.MessageHeader.Log()
// 	return append(ls, log.String("status", r.Status.String()))
// }
