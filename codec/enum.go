package codec

type SmsProto byte

const (
	CMPP20 SmsProto = iota
	CMPP21
	CMPP30
	SMGP13
	SMGP20
	SMGP30
	SGIP
	SMPP33
	SMPP34
)
