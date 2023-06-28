package sgip

import (
	"github.com/zhiyin2021/zysms/proto"
)

type SgipUnbindReq struct {
	SeqId []uint32
}
type SgipUnbindRsp struct {
	SeqId []uint32
}

func (u *SgipUnbindReq) Pack(seqId []uint32) []byte {
	data := make([]byte, SGIP_HEADER_LEN)
	pkt := proto.NewPacket(data)
	pkt.WriteU32(SGIP_HEADER_LEN)
	pkt.WriteU32(SGIP_UNBIND.ToInt())
	pkt.WriteU32(seqId[0])
	pkt.WriteU32(seqId[1])
	pkt.WriteU32(seqId[2])
	u.SeqId = seqId
	return data
}

func (u *SgipUnbindReq) Unpack(data []byte) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = r.(error)
		}
	}()
	var pkt = proto.NewPacket(data)
	// Sequence Id
	u.SeqId = make([]uint32, 3)
	u.SeqId[0] = pkt.ReadU32()
	u.SeqId[1] = pkt.ReadU32()
	u.SeqId[2] = pkt.ReadU32()
}

func (r *SgipUnbindRsp) Pack(seqId []uint32) []byte {
	data := make([]byte, SGIP_HEADER_LEN)
	pkt := proto.NewPacket(data)
	pkt.WriteU32(SGIP_HEADER_LEN)
	pkt.WriteU32(SGIP_UNBIND_RESP.ToInt())
	pkt.WriteU32(seqId[0])
	pkt.WriteU32(seqId[1])
	pkt.WriteU32(seqId[2])
	r.SeqId = seqId
	return data
}

func (r *SgipUnbindRsp) Unpack(data []byte) (e error) {
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
}
