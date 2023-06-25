package codec

type Version uint8

const (
	CMPP_HEADER_LEN  uint32 = 12
	CMPP2_PACKET_MAX uint32 = 2477
	CMPP2_PACKET_MIN uint32 = 12
	CMPP3_PACKET_MAX uint32 = 3335
	CMPP3_PACKET_MIN uint32 = 12

	V30 Version = 0x30
	V21 Version = 0x21
	V20 Version = 0x20
)

var versionStr = map[Version]string{
	V30: "cmpp30",
	V21: "cmpp21",
	V20: "cmpp20",
}

func (t Version) String() string {
	if v, ok := versionStr[t]; ok {
		return v
	}
	return "unknown"
}

// MajorMatch 主版本相匹配
func (t Version) MajorMatch(v uint8) bool {
	return uint8(t)&0xf0 == v&0xf0
}

// MajorMatchV 主版本相匹配
func (t Version) MajorMatchV(v Version) bool {
	return uint8(t)&0xf0 == uint8(v)&0xf0
}

// CommandId 命令定义
type CommandId uint32

const (
	CMPP_REQUEST_MIN, CMPP_RESPONSE_MIN CommandId = iota, 0x80000000 + iota
	CMPP_CONNECT, CMPP_CONNECT_RESP
	CMPP_TERMINATE, CMPP_TERMINATE_RESP
	_, _
	CMPP_SUBMIT, CMPP_SUBMIT_RESP
	CMPP_DELIVER, CMPP_DELIVER_RESP
	CMPP_QUERY, CMPP_QUERY_RESP
	CMPP_CANCEL, CMPP_CANCEL_RESP
	CMPP_ACTIVE_TEST, CMPP_ACTIVE_TEST_RESP
	CMPP_FWD, CMPP_FWD_RESP
	CMPP_REQUEST_MAX, CMPP_RESPONSE_MAX
)

var cmdStr = map[CommandId]string{
	CMPP_REQUEST_MIN:      "CMPP_REQUEST_MIN",
	CMPP_RESPONSE_MIN:     "CMPP_RESPONSE_MIN",
	CMPP_CONNECT:          "CMPP_CONNECT",
	CMPP_CONNECT_RESP:     "CMPP_CONNECT_RESP",
	CMPP_TERMINATE:        "CMPP_TERMINATE",
	CMPP_TERMINATE_RESP:   "CMPP_TERMINATE_RESP",
	CMPP_SUBMIT:           "CMPP_SUBMIT",
	CMPP_SUBMIT_RESP:      "CMPP_SUBMIT_RESP",
	CMPP_DELIVER:          "CMPP_DELIVER",
	CMPP_DELIVER_RESP:     "CMPP_DELIVER_RESP",
	CMPP_QUERY:            "CMPP_QUERY",
	CMPP_QUERY_RESP:       "CMPP_QUERY_RESP",
	CMPP_CANCEL:           "CMPP_CANCEL",
	CMPP_CANCEL_RESP:      "CMPP_CANCEL_RESP",
	CMPP_ACTIVE_TEST:      "CMPP_ACTIVE_TEST",
	CMPP_ACTIVE_TEST_RESP: "CMPP_ACTIVE_TEST_RESP",
	CMPP_FWD:              "CMPP_FWD",
	CMPP_FWD_RESP:         "CMPP_FWD_RESP",
	CMPP_REQUEST_MAX:      "CMPP_REQUEST_MAX",
	CMPP_RESPONSE_MAX:     "CMPP_RESPONSE_MAX",
}

func (id CommandId) ToInt() uint32 {
	return uint32(id)
}

func (id CommandId) String() string {
	if v, ok := cmdStr[id]; ok {
		return v
	}
	return "UNKNOWN"
}
