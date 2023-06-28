package sgip

import (
	"github.com/zhiyin2021/zysms/proto"
	"github.com/zhiyin2021/zysms/utils"
)

type SgipSubmitReq struct {
	SeqId            []uint32
	SPNumber         string   //  SP的接入号码【 21 bytes 】
	ChargeNumber     string   //  付费号码，手机号码前加“86”国别标志；当且仅当群发且对用户收费时为空；如果为空，则该条短消息产生的费用由UserNumber代表的用户支付；如果为全零字符串“000000000000000000000”，表示该条短消息产生的费用由SP支付。【 21 bytes 】
	UserCount        byte     //  接收短消息的手机数量，取值范围1至100【 1  bytes 】
	UserNumber       []string //  接收该短消息的手机号，该字段重复UserCount指定的次数，手机号码前加“86”国别标志【 21 bytes 】
	CorpId           string   //  企业代码，取值范围0-99999 【 5  bytes 】
	ServiceType      string   //  业务代码，由SP定义 【 10 bytes 】
	FeeType          byte     //  计费类型【 1  bytes 】
	FeeValue         string   //  取值范围0-99999，该条短消息的收费值，单位为分，由SP定义,对于包月制收费的用户，该值为月租费的值 【 6  bytes 】
	GivenValue       string   //  取值范围0-99999，赠送用户的话费，单位为分，由SP定义，特指由SP向用户发送广告时的赠送话费【 6  bytes 】
	AgentFlag        byte     //  代收费标志，0：应收；1：实收 【 1  bytes 】
	MorelatetoMTFlag byte     //  引起MT消息的原因 【 1  bytes 】
	Priority         byte     //  优先级0-9从低到高，默认为0 【 1 bytes 】
	ExpireTime       string   //  短消息寿命的终止时间，如果为空，表示使用短消息中心的缺省值。时间内容为16个字符，格式为”yymmddhhmmsstnnp” ，其中“tnnp”取固定值“032+”，即默认系统为北京时间 【 16 bytes 】
	ScheduleTime     string   //  短消息定时发送的时间，如果为空，表示立刻发送该短消息。时间内容为16个字符，格式为“yymmddhhmmsstnnp” ，其中“tnnp”取固定值“032+”，即默认系统为北京时间 【 16  bytes 】
	ReportFlag       byte     //  状态报告标记【 1  bytes 】
	TpPid            byte     //  GSM协议类型。详细解释请参考GSM03.40中的9.2.3.9 【 1  bytes 】
	TpUdhi           byte     //  GSM协议类型。详细解释请参考GSM03.40中的9.2.3.9 【 1  bytes 】
	MessageCoding    byte     //  短消息的编码格式。 【 1  bytes 】
	MessageType      byte     //  信息类型: 0-短消息信息 其它:待定 【 1  bytes 】
	MessageLength    uint32   //  短消息的长度【 4  bytes 】
	MessageContent   []byte   //  编码后消息内容
	Reserve          string   //  保留，扩展用【 8  bytes 】

	// ReportFlag
	// 状态报告标记 0-该条消息只有最后出错时要返回状态报告 1-该条消息无论最后是否成功都要返回状态报告 2-该条消息不需要返回状态报告 3-该条消息仅携带包月计费信息，不下发给用户， 要返回状态报告
	// 其它-保留
	// 缺省设置为 0

	// MorelatetoMTFlag
	// 引起 MT 消息的原因
	// 0-MO 点播引起的第一条 MT 消息;
	// 1-MO 点播引起的非第一条 MT 消息;
	// 2-非 MO 点播引起的 MT 消息;
	// 3-系统反馈引起的 MT 消息。

	// MessageCoding
	// 短消息的编码格式。
	// 0:纯 ASCII 字符串
	// 3:写卡操作
	// 4:二进制编码
	// 8:UCS2 编码
	// 15: GBK 编码
	// 其它参见 GSM3.38 第 4 节:SMS Data Coding Scheme
}

const (
	SgipSubmitReqLen = SGIP_HEADER_LEN + 123 // SGIP_HEADER_LEN + 123 + 21*len(s.UserNumber)+ MessageLength
	SgipSubmitRspLen = SGIP_HEADER_LEN + 9
	MtBaseLen        = 143
)

// func NewSubmit(ac *proto.AuthConf, phones []string, content string, options ...proto.OptionFunc) (messages []proto.RequestPdu) {
// 	mt := &Submit{}
// 	mt.PacketLength = MtBaseLen
// 	mt.CommandId = SGIP_SUBMIT
// 	mt.SequenceNumber = Sequencer.NextVal()
// 	mt.SetOptions(ac, proto.LoadMtOptions(options...))
// 	mt.UserCount = byte(len(phones))
// 	mt.UserNumber = phones
// 	mt.MessageCoding = utils.MsgFmt(content)
// 	mt.MorelatetoMTFlag = 2

// 	slices := utils.MsgSlices(mt.MessageCoding, content)
// 	if len(slices) == 1 {
// 		mt.MessageLength = uint32(len(slices[0]))
// 		mt.MessageContent = slices[0]
// 		mt.PacketLength = uint32(MtBaseLen + len(phones)*21 + len(slices[0]))
// 		return []proto.RequestPdu{mt}
// 	} else {
// 		mt.TpUdhi = 1
// 		for i, msgBytes := range slices {
// 			// 拷贝 mt
// 			tmp := *mt
// 			sub := &tmp
// 			if i != 0 {
// 				sub.SequenceNumber = Sequencer.NextVal()
// 			}
// 			sub.MessageLength = uint32(len(msgBytes))
// 			sub.MessageContent = msgBytes
// 			sub.PacketLength = uint32(MtBaseLen + len(phones)*21 + len(msgBytes))
// 			messages = append(messages, sub)
// 		}
// 		return messages
// 	}
// }

func (s *SgipSubmitReq) Pack(seqId []uint32) []byte {
	pktLen := int(SgipSubmitReqLen) + 21*len(s.UserNumber) + int(s.MessageLength)
	data := make([]byte, pktLen)
	pkt := proto.NewPacket(data)
	pkt.WriteU32(uint32(pktLen))
	pkt.WriteU32(SGIP_SUBMIT.ToInt())
	pkt.WriteU32(seqId[0])
	pkt.WriteU32(seqId[1])
	pkt.WriteU32(seqId[2])
	pkt.WriteStr(s.SPNumber, 21)
	pkt.WriteStr(s.ChargeNumber, 21)
	pkt.WriteByte(s.UserCount)
	for _, p := range s.UserNumber {
		pkt.WriteStr(p, 21)
	}
	pkt.WriteStr(s.CorpId, 5)
	pkt.WriteStr(s.ServiceType, 10)
	pkt.WriteByte(s.FeeType)
	pkt.WriteStr(s.FeeValue, 6)
	pkt.WriteStr(s.GivenValue, 6)
	pkt.WriteByte(s.AgentFlag)
	pkt.WriteByte(s.MorelatetoMTFlag)
	pkt.WriteByte(s.Priority)
	pkt.WriteStr(s.ExpireTime, 16)
	pkt.WriteStr(s.ScheduleTime, 16)
	pkt.WriteByte(s.ReportFlag)
	pkt.WriteByte(s.TpPid)
	pkt.WriteByte(s.TpUdhi)
	pkt.WriteByte(s.MessageCoding)
	pkt.WriteByte(s.MessageType)
	pkt.WriteU32(s.MessageLength)
	pkt.WriteBytes(s.MessageContent, int(s.MessageLength))
	pkt.WriteStr(s.Reserve, 8)
	return data
}

func (s *SgipSubmitReq) Unpack(data []byte) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = r.(error)
		}
	}()
	var pkt = proto.NewPacket(data)
	// Sequence Id
	s.SeqId = make([]uint32, 3)
	s.SeqId[0] = pkt.ReadU32()
	s.SeqId[1] = pkt.ReadU32()
	s.SeqId[2] = pkt.ReadU32()
	s.SPNumber = pkt.ReadStr(21)
	s.ChargeNumber = pkt.ReadStr(21)
	s.UserCount = pkt.ReadByte()
	s.UserNumber = make([]string, s.UserCount)
	for i := 0; i < int(s.UserCount); i++ {
		s.UserNumber[i] = pkt.ReadStr(21)
	}
	s.CorpId = pkt.ReadStr(5)
	s.ServiceType = pkt.ReadStr(10)
	s.FeeType = pkt.ReadByte()
	s.FeeValue = pkt.ReadStr(6)
	s.GivenValue = pkt.ReadStr(6)
	s.AgentFlag = pkt.ReadByte()
	s.MorelatetoMTFlag = pkt.ReadByte()
	s.Priority = pkt.ReadByte()
	s.ExpireTime = pkt.ReadStr(16)

	s.ScheduleTime = pkt.ReadStr(16)
	s.ReportFlag = pkt.ReadByte()
	s.TpPid = pkt.ReadByte()
	s.TpUdhi = pkt.ReadByte()
	s.MessageCoding = pkt.ReadByte()
	s.MessageType = pkt.ReadByte()
	s.MessageLength = pkt.ReadU32()
	content := pkt.ReadBytes(int(s.MessageLength))
	s.Reserve = pkt.ReadStr(8)
	s.MessageContent = content
	if content[0] == 0x05 && content[1] == 0x00 && content[2] == 0x03 {
		content = content[6:]
		s.MessageContent, _ = utils.Ucs2ToUtf8(content)
	}
	s.Reserve = ""
}

// func (s *Submit) SetOptions(ac *proto.AuthConf, ops *proto.MtOptions) {
// 	s.SPNumber = ac.SmsDisplayNo
// 	if ops.SpSubNo != "" {
// 		s.SPNumber += ops.SpSubNo
// 	}

// 	if len(ac.ClientId) > 5 {
// 		s.CorpId = ac.ClientId[5:]
// 	} else {
// 		s.CorpId = ac.ClientId
// 	}

// 	if ops.MsgLevel != uint8(0xf) {
// 		s.Priority = ops.MsgLevel
// 	} else {
// 		s.Priority = ac.DefaultMsgLevel
// 	}

// 	if ops.NeedReport != uint8(0xf) {
// 		s.ReportFlag = ops.NeedReport
// 	} else {
// 		s.ReportFlag = ac.NeedReport
// 	}

// 	s.ServiceType = ac.ServiceId
// 	if ops.ServiceId != "" {
// 		s.ServiceType = ops.ServiceId
// 	}

// 	if ops.AtTime != "" {
// 		s.ScheduleTime = ops.AtTime
// 	}

// 	if ops.ValidTime != "" {
// 		s.ExpireTime = ops.ValidTime
// 	} else {
// 		t := time.Now().Add(ac.MtValidDuration)
// 		s.ExpireTime = utils.FormatTime(t)
// 	}
// }

type SgipSubmitRsp struct {
	SeqId   []uint32 // 消息流水号，由SP接入的短信网关本身产生  12 bytes 】
	Status  Status   // 0：正确返回
	Reserve string   // 保留，扩展用【 8 bytes 】
}

func (r *SgipSubmitRsp) Unpack(data []byte) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = r.(error)
		}
	}()
	var pkt = proto.NewPacket(data)
	// Sequence Id
	r.SeqId = make([]uint32, 3)
	r.SeqId[0] = pkt.ReadU32()
	r.SeqId[1] = pkt.ReadU32()
	r.SeqId[2] = pkt.ReadU32()
	r.Status = Status(pkt.ReadByte())
	r.Reserve = pkt.ReadStr(8)
}

func (r *SgipSubmitRsp) Pack(seqId []uint32) []byte {
	data := make([]byte, SgipSubmitRspLen)
	pkt := proto.NewPacket(data)
	pkt.WriteU32(SgipSubmitRspLen)
	pkt.WriteU32(SGIP_SUBMIT_RESP.ToInt())
	pkt.WriteU32(seqId[0])
	pkt.WriteU32(seqId[1])
	pkt.WriteU32(seqId[2])
	pkt.WriteByte(byte(r.Status))
	pkt.WriteStr(r.Reserve, 8)
	return data
}
