package codec

import (
	"strconv"
	"time"
)

type msgId struct {
	tm       time.Time
	gateway  uint32
	sequence uint32
}

func NewMsgId(gatewayId uint32, seqId uint32) *msgId {
	return &msgId{
		tm:       time.Now(),
		gateway:  gatewayId,
		sequence: seqId,
	}
}

func (m *msgId) UInt64() uint64 {
	uid := uint64(0)
	uid |= uint64(m.tm.Month()) << 60
	uid |= uint64(m.tm.Day()) << 55
	uid |= uint64(m.tm.Hour()) << 50
	uid |= uint64(m.tm.Minute()) << 44
	uid |= uint64(m.tm.Second()) << 38
	uid |= uint64(m.gateway) << 16
	uid |= uint64(m.sequence) & 0xffff
	return uint64(uid)
}

func (m *msgId) String() string {
	return strconv.FormatUint(m.UInt64(), 16)
}
