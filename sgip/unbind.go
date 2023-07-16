package sgip

import "github.com/zhiyin2021/zysms/codec"

type SgipUnbindReq struct {
	SeqId []uint32
}
type SgipUnbindRsp struct {
	SeqId []uint32
}

func (u *SgipUnbindReq) Pack(seqId uint32) []byte {
	pkt := codec.NewWriter(SGIP_HEADER_LEN, SGIP_UNBIND.ToInt())
	pkt.WriteU32(seqId)
	pkt.WriteU32(getTm())
	pkt.WriteU32(seqId)
	return pkt.Bytes()
}

func (u *SgipUnbindReq) Unpack(data []byte) error {
	var pkt = codec.NewReader(data)
	// Sequence Id
	u.SeqId = make([]uint32, 3)
	u.SeqId[0] = pkt.ReadU32()
	u.SeqId[1] = pkt.ReadU32()
	u.SeqId[2] = pkt.ReadU32()
	return pkt.Err()
}

func (r *SgipUnbindRsp) Pack(seqId uint32) []byte {
	pkt := codec.NewWriter(SGIP_HEADER_LEN, SGIP_UNBIND_RESP.ToInt())
	pkt.WriteU32(seqId)

	pkt.WriteU32(getTm())
	pkt.WriteU32(seqId)
	return pkt.Bytes()
}

func (r *SgipUnbindRsp) Unpack(data []byte) error {
	var pkt = codec.NewReader(data)
	// Sequence Id
	r.SeqId = make([]uint32, 3)
	r.SeqId[0] = pkt.ReadU32()
	r.SeqId[1] = pkt.ReadU32()
	r.SeqId[2] = pkt.ReadU32()
	return pkt.Err()
}
