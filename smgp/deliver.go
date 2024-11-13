package smgp

import (
	"fmt"
	"strings"

	"github.com/zhiyin2021/zysms/codec"
)

type DeliverReq struct {
	base
	MsgId      string // 【10字节】短消息流水号
	IsReport   byte   // 【1字节】是否为状态报告
	MsgFormat  byte   // 【1字节】短消息格式
	RecvTime   string // 【14字节】短消息定时发送时间
	SrcTermID  string // 【21字节】短信息发送方号码
	DestTermID string // 【21】短消息接收号码

	Message codec.ShortMessage // 【MsgLength字节】短消息内容
	// MsgBytes   []byte         // 消息内容按照Msg_Fmt编码后的数据
	//Report  *Report        // 状态报告
	Reserve string // 【8字节】保留

	Report *DeliverReport

	// 协议版本,不是报文内容，但在调用encode方法前需要设置此值
	// Version Version
}

type DeliverResp struct {
	base
	MsgId  string // 【10字节】短消息流水号
	Status Status

	// 协议版本,不是报文内容，但在调用encode方法前需要设置此值
	// Version Version
}

func NewDeliverReq(ver codec.Version) codec.PDU {
	return &DeliverReq{
		base: newBase(ver, SMGP_DELIVER, 0),
	}
}
func NewDeliverResp(ver codec.Version) codec.PDU {
	return &DeliverResp{
		base: newBase(ver, SMGP_DELIVER_RESP, 0),
	}
}

// Pack packs the ActiveTestReq to bytes stream for client side.
func (p *DeliverReq) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteStr(p.MsgId, 10)
		if p.Report != nil {
			p.IsReport = 1
			p.encodeReport()
		}
		bw.WriteByte(p.IsReport)
		bw.WriteByte(p.MsgFormat)
		bw.WriteStr(p.RecvTime, 14)
		bw.WriteStr(p.SrcTermID, 21)
		bw.WriteStr(p.DestTermID, 21)
		p.Message.Marshal(bw)

		bw.WriteStr(p.Reserve, 8)
	})
}

// Unpack unpack the binary byte stream to a ActiveTestReq variable.
// After unpack, you will get all value of fields in
// ActiveTestReq struct.
func (p *DeliverReq) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.MsgId = br.ReadStr(10)
		p.IsReport = br.ReadU8()
		p.MsgFormat = br.ReadU8()
		p.RecvTime = br.ReadStr(14)
		p.SrcTermID = br.ReadStr(21)
		p.DestTermID = br.ReadStr(21)
		p.Message.Unmarshal(br, false, p.MsgFormat)
		if p.IsReport == 1 {
			p.decodeReport()
		}
		p.Reserve = br.ReadStr(8)
		return br.Err()
	})
}

// GetResponse implements PDU interface.
func (b *DeliverReq) GetResponse() codec.PDU {
	return &DeliverResp{
		base: newBase(b.Version, SMGP_DELIVER_RESP, b.SequenceNumber),
	}
}

// Pack packs the ActiveTestReq to bytes stream for client side.
func (p *DeliverResp) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteStr(p.MsgId, 10)
		bw.WriteU32(uint32(p.Status))
	})
}

// Unpack unpack the binary byte stream to a ActiveTestReq variable.
// After unpack, you will get all value of fields in
// ActiveTestReq struct.
func (p *DeliverResp) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.MsgId = br.ReadStr(10)
		p.Status = Status(br.ReadU32())
		return br.Err()
	})
}

// GetResponse implements PDU interface.
func (b *DeliverResp) GetResponse() codec.PDU {
	return nil
}

type DeliverReport struct {
	MsgId      string // 10字节 The message ID allocated to the message by the SMSC when originally submitted.
	Sub        string // 3 字节 Number of short messages originally submitted. This is only relevant when the original message was submitted to a distribution list.The value is padded with leading zeros if necessary.
	Dlvrd      string // 3 字节 Number of short messages delivered. This is only relevant where the original message was submitted to a distribution list.The value is padded with leading zeros if necessary.
	SubmitDate string // 10字节 (YYMMDDhhmm)The time and date at which the short message was submitted. In the case of a message which has been replaced, this is the date that the original message was replaced.The format is as follows:
	DoneDate   string // 10字节 (YYMMDDhhmm)The time and date at which the short message reached it’s final state. The format is the same as for the submit date.
	Stat       string // 7 字节 The final status of the message. For settings for this field see Table B-2.
	Err        string // 3 字节
	Text       string // 20字节 The first 20 characters of the short message.
}

func (c *DeliverReq) decodeReport() {
	c.Report = &DeliverReport{}
	msg := c.Message.GetMessage()
	c.Report.MsgId, msg = splitReport(msg, "id:")
	c.Report.Sub, msg = splitReport(msg, "sub:")
	c.Report.Dlvrd, msg = splitReport(msg, "dlvrd:")
	c.Report.SubmitDate, msg = splitReport(msg, "submit date:")
	c.Report.DoneDate, msg = splitReport(msg, "done date:")
	c.Report.Stat, msg = splitReport(msg, "stat:")
	c.Report.Err, msg = splitReport(msg, "err:")
	c.Report.Text, _ = splitReport(msg, "text:")
}
func (c *DeliverReq) encodeReport() {
	if c.Report != nil {
		//fmt.Sprintf("id:%s sub:%s dlvrd:%s submit date:%s done date:%s stat:%s err:%s text:%s ", c.Report.MsgId, c.Report.Sub, c.Report.Dlvrd, c.Report.SubmitDate, c.Report.DoneDate, c.Report.Stat, c.Report.Text)
		c.Message.SetMessage(c.Report.String(), codec.ASCII)
	}
}

func splitReport(content, sub1 string) (retSub string, retContent string) {
	n := strings.Index(content, sub1)
	if n == -1 {
		return content, ""
	}
	n += len(sub1)
	m := strings.Index(content[n:], " ")
	if m == -1 {
		return content, ""
	}
	return content[n : m+n], content[n+m:]
}
func (r *DeliverReport) String() string {
	if r.Err == "" {
		r.Err = "0"
	}
	return fmt.Sprintf("id:%s sub:%s dlvrd:%s submit date:%s done date:%s stat:%s err:%s text:%s ", r.MsgId, r.Sub, r.Dlvrd, r.SubmitDate, r.DoneDate, r.Stat, r.Err, r.Text)
}
