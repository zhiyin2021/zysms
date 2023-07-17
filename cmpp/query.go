package cmpp

import (
	"github.com/zhiyin2021/zysms/codec"
)

type QueryReq struct {
	base
	Time      string //	8字节 YYYYMMDD
	QueryType byte   //	1字节 0：总数查询；1：按业务类型查询
	QueryCode string // 10字节  查询码 当 QueryType 为 0 时，此项无效;当 QueryType 为 1 时，此项填写业务类 型 Service_Id.
	Reserve   string // 8字节 保留
}
type QueryResp struct {
	base
	Time      string //	8字节 YYYYMMDD
	QueryType byte   //	1字节 0：总数查询；1：按业务类型查询
	QueryCode string // 10字节  查询码 当 QueryType 为 0 时，此项无效;当 QueryType 为 1 时，此项填写业务类 型 Service_Id.

	MtTlMsg uint32 // 4字节 从SP接收信息总数
	MtTlUsr uint32 // 4字节 从SP接收用户总数
	MtScs   uint32 // 4字节 成功转发数量
	MtWt    uint32 // 4字节 待转发数量
	MtFl    uint32 // 4字节 转发失败数量
	MoScs   uint32 // 4字节 向SP成功送达数量
	MoWt    uint32 // 4字节 向SP待送达数量
	MoFl    uint32 // 4字节 向SP送达失败数量

}

func NewQueryReq(ver Version) codec.PDU {
	return &QueryReq{
		base: newBase(ver, CMPP_QUERY, 0),
	}
}
func NewQueryResp(ver Version) codec.PDU {
	return &QueryResp{
		base: newBase(ver, CMPP_QUERY_RESP, 0),
	}
}

// Pack packs the ActiveTestReq to bytes stream for client side.
func (p *QueryReq) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteStr(p.Time, 8)
		bw.WriteByte(p.QueryType)
		bw.WriteStr(p.QueryCode, 10)
		bw.WriteStr(p.Reserve, 8)
	})
}

// Unpack unpack the binary byte stream to a ActiveTestReq variable.
// After unpack, you will get all value of fields in
// ActiveTestReq struct.
func (p *QueryReq) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.Time = br.ReadStr(8)
		p.QueryType = br.ReadByte()
		p.QueryCode = br.ReadStr(10)
		p.Reserve = br.ReadStr(8)
		return br.Err()
	})
}
func (p *QueryReq) GetResponse() codec.PDU {
	return &QueryResp{
		base: newBase(p.Version, CMPP_QUERY_RESP, p.SequenceNumber),
	}
}

// Pack packs the ActiveTestReq to bytes stream for client side.
func (p *QueryResp) Marshal(w *codec.BytesWriter) {
	p.base.marshal(w, func(bw *codec.BytesWriter) {
		bw.WriteStr(p.Time, 8)
		bw.WriteByte(p.QueryType)
		bw.WriteStr(p.QueryCode, 10)
		bw.WriteU32(p.MtTlMsg)
		bw.WriteU32(p.MtTlUsr)
		bw.WriteU32(p.MtScs)
		bw.WriteU32(p.MtWt)
		bw.WriteU32(p.MtFl)
		bw.WriteU32(p.MoScs)
		bw.WriteU32(p.MoWt)
		bw.WriteU32(p.MoFl)
	})
}

// Unpack unpack the binary byte stream to a ActiveTestReq variable.
// After unpack, you will get all value of fields in
// ActiveTestReq struct.
func (p *QueryResp) Unmarshal(w *codec.BytesReader) error {
	return p.base.unmarshal(w, func(br *codec.BytesReader) error {
		p.Time = br.ReadStr(8)
		p.QueryType = br.ReadByte()
		p.QueryCode = br.ReadStr(10)
		p.MtTlMsg = br.ReadU32()
		p.MtTlUsr = br.ReadU32()
		p.MtScs = br.ReadU32()
		p.MtWt = br.ReadU32()
		p.MtFl = br.ReadU32()
		p.MoScs = br.ReadU32()
		p.MoWt = br.ReadU32()
		p.MoFl = br.ReadU32()
		return br.Err()
	})
}
func (p *QueryResp) GetResponse() codec.PDU {
	return nil
}
