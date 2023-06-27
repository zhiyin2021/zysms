package cmpp

import (
	"github.com/zhiyin2021/zysms/event"
	"github.com/zhiyin2021/zysms/proto"
)

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

func (id CommandId) ToInt() uint32 {
	return uint32(id)
}

func (id CommandId) String() string {
	if v, ok := cmdStr[id]; ok {
		return v
	}
	return "UNKNOWN"
}
func (id CommandId) Event() event.SmsEvent {
	if v, ok := cmdEvent[id]; ok {
		return v
	}
	return event.SmsEventUnknown
}

var cmdEvent = map[CommandId]event.SmsEvent{
	CMPP_CONNECT:          event.SmsEventAuthReq,       //   "CMPP_CONNECT",
	CMPP_CONNECT_RESP:     event.SmsEventAuthRsp,       // "CMPP_CONNECT_RESP",
	CMPP_TERMINATE:        event.SmsEventTerminateReq,  //   "CMPP_TERMINATE",
	CMPP_TERMINATE_RESP:   event.SmsEventActiveTestRsp, //   "CMPP_TERMINATE_RESP",
	CMPP_SUBMIT:           event.SmsEventSubmitReq,     //   "CMPP_SUBMIT",
	CMPP_SUBMIT_RESP:      event.SmsEventSubmitRsp,     //   "CMPP_SUBMIT_RESP",
	CMPP_DELIVER:          event.SmsEventDeliverReq,    //   "CMPP_DELIVER",
	CMPP_DELIVER_RESP:     event.SmsEventDeliverRsp,    //   "CMPP_DELIVER_RESP",
	CMPP_QUERY:            event.SmsEventQueryReq,      //   "CMPP_QUERY",
	CMPP_QUERY_RESP:       event.SmsEventQueryRsp,      //   "CMPP_QUERY_RESP",
	CMPP_CANCEL:           event.SmsEventCancelReq,     //   "CMPP_CANCEL",
	CMPP_CANCEL_RESP:      event.SmsEventCancelRsp,     //   "CMPP_CANCEL_RESP",
	CMPP_ACTIVE_TEST:      event.SmsEventActiveTestReq, //   "CMPP_ACTIVE_TEST",
	CMPP_ACTIVE_TEST_RESP: event.SmsEventActiveTestRsp, //  "CMPP_ACTIVE_TEST_RESP",
}
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
var CmppPacket = map[CommandId]func(Version, []byte) proto.Packer{
	CMPP_REQUEST_MIN:      nil,                  //"CMPP_REQUEST_MIN",
	CMPP_RESPONSE_MIN:     nil,                  //    "CMPP_RESPONSE_MIN",
	CMPP_CONNECT:          newCmppConnReq,       //   "CMPP_CONNECT",
	CMPP_CONNECT_RESP:     newCmppConnRsp,       // "CMPP_CONNECT_RESP",
	CMPP_TERMINATE:        newCmppTerminateReq,  //   "CMPP_TERMINATE",
	CMPP_TERMINATE_RESP:   newCmppTerminateRsp,  //   "CMPP_TERMINATE_RESP",
	CMPP_SUBMIT:           newCmppSubmitReq,     //   "CMPP_SUBMIT",
	CMPP_SUBMIT_RESP:      newCmppSubmitRsp,     //   "CMPP_SUBMIT_RESP",
	CMPP_DELIVER:          newCmppDeliverReq,    //   "CMPP_DELIVER",
	CMPP_DELIVER_RESP:     newCmppDeliverRsp,    //   "CMPP_DELIVER_RESP",
	CMPP_QUERY:            newCmppQueryReq,      //   "CMPP_QUERY",
	CMPP_QUERY_RESP:       newCmppQueryRsp,      //   "CMPP_QUERY_RESP",
	CMPP_CANCEL:           newCmppCancelReq,     //   "CMPP_CANCEL",
	CMPP_CANCEL_RESP:      newCmppCancelRsp,     //   "CMPP_CANCEL_RESP",
	CMPP_ACTIVE_TEST:      newCmppActiveTestReq, //   "CMPP_ACTIVE_TEST",
	CMPP_ACTIVE_TEST_RESP: newCmppActiveTestRsp, //  "CMPP_ACTIVE_TEST_RESP",
	CMPP_FWD:              newCmppFwdReq,        //  "CMPP_FWD",
	CMPP_FWD_RESP:         newCmppFwdRsp,        //   "CMPP_FWD_RESP",
	CMPP_REQUEST_MAX:      nil,                  //   "CMPP_REQUEST_MAX",
	CMPP_RESPONSE_MAX:     nil,                  //  "CMPP_RESPONSE_MAX",
}

func newCmppConnReq(v Version, data []byte) proto.Packer {
	return (&CmppConnReq{}).Unpack(data)
}
func newCmppConnRsp(v Version, data []byte) proto.Packer {
	if v == V30 {
		return (&Cmpp3ConnRsp{}).Unpack(data)
	}
	return (&Cmpp2ConnRsp{}).Unpack(data)
}
func newCmppTerminateReq(v Version, data []byte) proto.Packer {
	return (&CmppConnReq{}).Unpack(data)
}
func newCmppTerminateRsp(v Version, data []byte) proto.Packer {
	return (&CmppTerminateRsp{}).Unpack(data)
}
func newCmppSubmitReq(v Version, data []byte) proto.Packer {
	if v == V30 {
		return (&Cmpp3SubmitReq{}).Unpack(data)
	}
	return (&Cmpp2SubmitReq{}).Unpack(data)
}
func newCmppSubmitRsp(v Version, data []byte) proto.Packer {
	if v == V30 {
		return (&Cmpp3SubmitRsp{}).Unpack(data)
	}
	return (&Cmpp2SubmitRsp{}).Unpack(data)
}
func newCmppDeliverReq(v Version, data []byte) proto.Packer {
	if v == V30 {
		return (&Cmpp3DeliverReq{}).Unpack(data)
	}
	return (&Cmpp2DeliverReq{}).Unpack(data)
}
func newCmppDeliverRsp(v Version, data []byte) proto.Packer {
	if v == V30 {
		return (&Cmpp3DeliverRsp{}).Unpack(data)
	}
	return (&Cmpp2DeliverRsp{}).Unpack(data)
}
func newCmppQueryReq(v Version, data []byte) proto.Packer {
	return (&CmppQueryReq{}).Unpack(data)
}
func newCmppQueryRsp(v Version, data []byte) proto.Packer {
	return (&CmppQueryRsp{}).Unpack(data)
}
func newCmppCancelReq(v Version, data []byte) proto.Packer {
	return (&CmppCancelReq{}).Unpack(data)
}
func newCmppCancelRsp(v Version, data []byte) proto.Packer {
	return (&CmppCancelRsp{}).Unpack(data)
}
func newCmppActiveTestReq(v Version, data []byte) proto.Packer {
	return (&CmppActiveTestReq{}).Unpack(data)
}
func newCmppActiveTestRsp(v Version, data []byte) proto.Packer {
	return (&CmppActiveTestRsp{}).Unpack(data)
}
func newCmppFwdReq(v Version, data []byte) proto.Packer {
	if v == V30 {
		return (&Cmpp3FwdReq{}).Unpack(data)
	}
	return (&Cmpp2FwdReq{}).Unpack(data)
}
func newCmppFwdRsp(v Version, data []byte) proto.Packer {
	if v == V30 {
		return (&Cmpp3FwdRsp{}).Unpack(data)
	}
	return (&Cmpp2FwdRsp{}).Unpack(data)
}
