package smgp

import "github.com/zhiyin2021/zysms/codec"

type SmgpExitReq struct {
	seqId uint32
}
type SmgpExitRsp struct {
	seqId uint32
}

func (p *SmgpExitReq) Pack(seqId uint32) []byte {
	pkt := codec.NewWriter(SMGP_HEADEER_LEN, SMGP_EXIT.ToInt())
	pkt.WriteU32(seqId)
	p.seqId = seqId
	return pkt.Bytes()
}

func (p *SmgpExitReq) Unpack(data []byte) error {
	var pkt = codec.NewReader(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	return pkt.Err()
}
func (p *SmgpExitReq) SeqId() uint32 {
	return p.seqId
}

func (p *SmgpExitRsp) Pack(seqId uint32) []byte {
	pkt := codec.NewWriter(SMGP_HEADEER_LEN, SMGP_EXIT_RESP.ToInt())
	pkt.WriteU32(seqId)
	p.seqId = seqId
	return pkt.Bytes()
}

func (p *SmgpExitRsp) Unpack(data []byte) error {
	var pkt = codec.NewReader(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	return pkt.Err()
}

func (p *SmgpExitRsp) SeqId() uint32 {
	return p.seqId
}
