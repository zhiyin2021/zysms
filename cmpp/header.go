package cmpp

import "encoding/binary"

type CmppHeader struct {
	TotalLength uint32    // 4 bytes
	CommandId   CommandId // 4 bytes
	SeqId       uint32    // 4 bytes
}

// ParseHeader parses PDU header.
func ParseHeader(v [16]byte) (h CmppHeader) {
	h.TotalLength = binary.BigEndian.Uint32(v[:])
	h.CommandId = CommandId(binary.BigEndian.Uint32(v[4:]))
	h.SeqId = binary.BigEndian.Uint32(v[8:])
	return
}
