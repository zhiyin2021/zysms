package smgp

import "github.com/zhiyin2021/zysms/proto"

type SmgpExitReq struct {
	seqId uint32
}
type SmgpExitRsp struct {
	seqId uint32
}

func (p *SmgpExitReq) Pack(seqId uint32) []byte {
	data := make([]byte, SMGP_HEADEER_LEN)
	pkt := proto.NewPacket(data)
	pkt.WriteU32(SMGP_HEADEER_LEN)
	pkt.WriteU32(SMGP_EXIT.ToInt())

	p.seqId = seqId

	pkt.WriteU32(p.seqId)
	return data
}

func (p *SmgpExitReq) Unpack(data []byte) (e error) {
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
func (p *SmgpExitReq) SeqId() uint32 {
	return p.seqId
}

func (p *SmgpExitRsp) Pack(seqId uint32) []byte {
	data := make([]byte, SMGP_HEADEER_LEN)
	pkt := proto.NewPacket(data)
	pkt.WriteU32(SMGP_HEADEER_LEN)
	pkt.WriteU32(SMGP_EXIT_RESP.ToInt())
	if seqId > 0 {
		p.seqId = seqId
	}
	pkt.WriteU32(p.seqId)

	return data
}

func (p *SmgpExitRsp) Unpack(data []byte) (e error) {
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

func (p *SmgpExitRsp) SeqId() uint32 {
	return p.seqId
}
