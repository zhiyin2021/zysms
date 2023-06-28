package smgp

import "github.com/zhiyin2021/zysms/proto"

type SmgpActiveTest struct {
	seqId uint32
}
type SmgpActiveTestRsp struct {
	seqId uint32
}

func (p *SmgpActiveTest) Pack(seqId uint32) []byte {
	data := make([]byte, SMGP_HEADEER_LEN)
	pkt := proto.NewPacket(data)
	pkt.WriteU32(SMGP_HEADEER_LEN)
	pkt.WriteU32(SMGP_ACTIVE_TEST.ToInt())
	if seqId > 0 {
		p.seqId = seqId
	}
	pkt.WriteU32(p.seqId)
	return data
}

func (p *SmgpActiveTest) Unpack(data []byte) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = r.(error)
		}
	}()
	var pkt = proto.NewPacket(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	return p
}
func (p *SmgpActiveTest) SeqId() uint32 {
	return p.seqId
}
func (p *SmgpActiveTestRsp) Pack(seqId uint32) []byte {
	data := make([]byte, SMGP_HEADEER_LEN)
	pkt := proto.NewPacket(data)
	pkt.WriteU32(SMGP_HEADEER_LEN)
	pkt.WriteU32(SMGP_ACTIVE_TEST_RESP.ToInt())

	p.seqId = seqId

	pkt.WriteU32(p.seqId)
	p.seqId = seqId
	return data
}

func (p *SmgpActiveTestRsp) Unpack(data []byte) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = r.(error)
		}
	}()
	var pkt = proto.NewPacket(data)
	// Sequence Id
	p.seqId = pkt.ReadU32()
	return p
}

func (p *SmgpActiveTestRsp) SeqId() uint32 {
	return p.seqId
}
