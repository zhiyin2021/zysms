package sgip

import (
	"github.com/zhiyin2021/zysms/codec"
)

type SubmitReq struct {
	base
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
	Message          codec.ShortMessage
	Reserve          string //  保留，扩展用【 8  bytes 】

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
type SubmitResp struct {
	base
	Status  Status // 0：正确返回
	Reserve string // 保留，扩展用【 8 bytes 】
}

func NewSubmitReq(ver codec.Version, nodeId uint32) codec.PDU {
	c := &SubmitReq{
		base: newBase(ver, SGIP_SUBMIT, [3]uint32{nodeId, 0, 0}),
	}
	return c
}
func NewSubmitResp(ver codec.Version, nodeId uint32) codec.PDU {
	c := &SubmitResp{
		base: newBase(ver, SGIP_SUBMIT_RESP, [3]uint32{nodeId, 0, 0}),
	}
	return c
}

func (s *SubmitReq) Marshal(w *codec.BytesWriter) {
	s.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteStr(s.SPNumber, 21)
		bw.WriteStr(s.ChargeNumber, 21)
		bw.WriteByte(s.UserCount)
		for _, p := range s.UserNumber {
			bw.WriteStr(p, 21)
		}
		bw.WriteStr(s.CorpId, 5)
		bw.WriteStr(s.ServiceType, 10)
		bw.WriteByte(s.FeeType)
		bw.WriteStr(s.FeeValue, 6)
		bw.WriteStr(s.GivenValue, 6)
		bw.WriteByte(s.AgentFlag)
		bw.WriteByte(s.MorelatetoMTFlag)
		bw.WriteByte(s.Priority)
		bw.WriteStr(s.ExpireTime, 16)
		bw.WriteStr(s.ScheduleTime, 16)
		bw.WriteByte(s.ReportFlag)
		bw.WriteByte(s.TpPid)
		bw.WriteByte(s.TpUdhi)
		bw.WriteByte(s.MessageCoding)
		bw.WriteByte(s.MessageType)
		s.Message.Marshal(bw)
		bw.WriteStr(s.Reserve, 8)
	})
}

func (p *SubmitReq) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.SPNumber = br.ReadStr(21)
		p.ChargeNumber = br.ReadStr(21)
		p.UserCount = br.ReadByte()
		p.UserNumber = make([]string, p.UserCount)
		for i := 0; i < int(p.UserCount); i++ {
			p.UserNumber[i] = br.ReadStr(21)
		}
		p.CorpId = br.ReadStr(5)
		p.ServiceType = br.ReadStr(10)
		p.FeeType = br.ReadByte()
		p.FeeValue = br.ReadStr(6)
		p.GivenValue = br.ReadStr(6)
		p.AgentFlag = br.ReadByte()
		p.MorelatetoMTFlag = br.ReadByte()
		p.Priority = br.ReadByte()
		p.ExpireTime = br.ReadStr(16)

		p.ScheduleTime = br.ReadStr(16)
		p.ReportFlag = br.ReadByte()
		p.TpPid = br.ReadByte()
		p.TpUdhi = br.ReadByte()
		p.MessageCoding = br.ReadByte()
		p.MessageType = br.ReadByte()
		p.Message.Unmarshal(br, p.TpUdhi == 1, p.MessageCoding)
		p.Reserve = br.ReadStr(8)
		return br.Err()
	})
}

func (b *SubmitReq) GetResponse() codec.PDU {
	return &SubmitResp{
		base: newBase(b.Version, SGIP_SUBMIT_RESP, b.SequenceNumber),
	}
}

func (p *SubmitResp) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteByte(byte(p.Status))
		bw.WriteStr(p.Reserve, 8)
	})
}

func (p *SubmitResp) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.Status = Status(br.ReadByte())
		p.Reserve = br.ReadStr(8)
		return br.Err()
	})
}
func (p *SubmitResp) GetResponse() codec.PDU {
	return nil
}

// func (s *Submit) SetOptions(ac *codec.AuthConf, ops *codec.MtOptions) {
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
