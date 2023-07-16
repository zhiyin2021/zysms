package smgp

import "github.com/zhiyin2021/zysms/codec"

type SmgpActiveTest struct {
	seqId uint32
}
type SmgpActiveTestRsp struct {
	seqId uint32
}

func (p *SmgpActiveTest) Pack(seqId uint32) []byte {
	pkt := codec.NewWriter(SMGP_HEADEER_LEN, SMGP_ACTIVE_TEST.ToInt())
	pkt.WriteU32(seqId)
	p.seqId = seqId
	return pkt.Bytes()
}

func (p *SmgpActiveTest) Unpack(data []byte) error {
	var pkt = codec.NewReader(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	return pkt.Err()
}
func (p *SmgpActiveTest) SeqId() uint32 {
	return p.seqId
}
func (p *SmgpActiveTestRsp) Pack(seqId uint32) []byte {
	pkt := codec.NewWriter(SMGP_HEADEER_LEN, SMGP_ACTIVE_TEST_RESP.ToInt())
	pkt.WriteU32(seqId)
	p.seqId = seqId
	return pkt.Bytes()
}

func (p *SmgpActiveTestRsp) Unpack(data []byte) error {
	var pkt = codec.NewReader(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	return pkt.Err()
}

func (p *SmgpActiveTestRsp) SeqId() uint32 {
	return p.seqId
}
