package cmpp

type CmppHeader struct {
	TotalLength uint32 // 4 bytes
	CommandId   uint32 // 4 bytes
	SequenceId  uint32 // 4 bytes
}
