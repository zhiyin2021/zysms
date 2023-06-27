package proto

type SmsProto byte

const (
	CMPP2 SmsProto = iota
	CMPP3
	SMGP
	SGIP
	SMPP
)
