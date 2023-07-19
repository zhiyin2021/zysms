package sgip

import "github.com/zhiyin2021/zysms/codec"

type ReportReq struct {
	base
	ReportType byte
	UserNumber string
	State      Status
	ErrorCode  byte
	Reserve    string
}
type ReportResp struct {
	base
	Status  Status
	Reserve string
}

func NewReportReq(ver codec.Version, nodeId uint32) codec.PDU {
	return &ReportReq{
		base: newBase(ver, SGIP_REPORT, [3]uint32{nodeId, 0, 0}),
	}
}
func NewReportResp(ver codec.Version, nodeId uint32) codec.PDU {
	return &ReportResp{
		base: newBase(ver, SGIP_REPORT_RESP, [3]uint32{nodeId, 0, 0}),
	}
}
func (p *ReportReq) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteByte(p.ReportType)
		bw.WriteStr(p.UserNumber, 21)
		bw.WriteByte(byte(p.State))
		bw.WriteByte(p.ErrorCode)
		bw.WriteStr(p.Reserve, 8)
	})
}

func (p *ReportReq) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.ReportType = br.ReadByte()
		p.UserNumber = br.ReadStr(21)
		p.State = Status(br.ReadByte())
		p.ErrorCode = br.ReadByte()
		p.Reserve = br.ReadStr(8)
		return br.Err()
	})
}

func (b *ReportReq) GetResponse() codec.PDU {
	return &ReportResp{
		base: newBase(b.Version, SGIP_REPORT_RESP, b.SequenceNumber),
	}
}

func (p *ReportResp) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteByte(byte(p.Status))
		bw.WriteStr(p.Reserve, 8)
	})
}

func (p *ReportResp) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.Status = Status(br.ReadByte())
		p.Reserve = br.ReadStr(8)
		return br.Err()
	})
}

func (b *ReportResp) GetResponse() codec.PDU {
	return nil
}
