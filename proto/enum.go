package proto

type SmsProto byte

const (
	CMPP20 SmsProto = iota
	CMPP21
	CMPP30
	SMGP
	SGIP
	SMPP
)
