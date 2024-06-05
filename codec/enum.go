package codec

type SmsProto byte

const (
	UNKNOWN SmsProto = iota
	CMPP20
	CMPP21
	CMPP30
	SMGP13
	SMGP20
	SMGP30
	SGIP
	SMPP33
	SMPP34
)

var protoMap = map[SmsProto]Version{
	CMPP20: 0x20,
	CMPP21: 0x21,
	CMPP30: 0x30,
	SMGP13: 0x13,
	SMGP20: 0x20,
	SMGP30: 0x30,
	SGIP:   0x12,
	SMPP33: 0x33,
	SMPP34: 0x34,
}
var protoVer = map[SmsProto]string{
	CMPP20: "cmpp20",
	CMPP21: "cmpp21",
	CMPP30: "cmpp30",
	SMGP13: "smgp13",
	SMGP20: "smgp20",
	SMGP30: "smgp30",
	SGIP:   "sgip",
	SMPP33: "smpp33",
	SMPP34: "smpp34",
}

var proto = map[SmsProto]string{
	CMPP20: "cmpp",
	CMPP21: "cmpp",
	CMPP30: "cmpp",
	SMGP13: "smgp",
	SMGP20: "smgp",
	SMGP30: "smgp",
	SGIP:   "sgip",
	SMPP33: "smpp",
	SMPP34: "smpp",
}

func (s SmsProto) Version() Version {
	if v, ok := protoMap[s]; ok {
		return v
	}
	return 0
}

func (s SmsProto) String() string {
	if v, ok := protoVer[s]; ok {
		return v
	}
	return "unknown"
}

func (s SmsProto) Raw() string {
	if v, ok := proto[s]; ok {
		return v
	}
	return "unknown"
}
