package cmpp

import (
	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/smserror"
)

const (
	// CMPP_HEADER_LEN  uint32 = 12
	CMPP2_PACKET_MAX uint32 = 2477
	CMPP3_PACKET_MAX uint32 = 3335

	V30 codec.Version = 0x30
	V21 codec.Version = 0x21
	V20 codec.Version = 0x20

	SM_MSG_LEN      = 140
	PDU_HEADER_SIZE = 12
	MAX_PDU_LEN     = 3335
)

var versionStr = map[codec.Version]string{
	V30: "cmpp30",
	V21: "cmpp21",
	V20: "cmpp20",
}

// CommandId 命令定义

const (
	CMPP_REQUEST_MIN, CMPP_RESPONSE_MIN codec.CommandId = iota, 0x80000000 + iota
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

var pduMap = map[codec.CommandId]func(codec.Version) codec.PDU{
	CMPP_REQUEST_MIN:      nil,               //"CMPP_REQUEST_MIN",
	CMPP_RESPONSE_MIN:     nil,               //    "CMPP_RESPONSE_MIN",
	CMPP_CONNECT:          NewConnReq,        //   "CMPP_CONNECT",
	CMPP_CONNECT_RESP:     NewConnResp,       // "CMPP_CONNECT_RESP",
	CMPP_TERMINATE:        NewTerminateReq,   //   "CMPP_TERMINATE",
	CMPP_TERMINATE_RESP:   NewTerminateResp,  //   "CMPP_TERMINATE_RESP",
	CMPP_SUBMIT:           NewSubmitReq,      //   "CMPP_SUBMIT",
	CMPP_SUBMIT_RESP:      NewSubmitResp,     //   "CMPP_SUBMIT_RESP",
	CMPP_DELIVER:          NewDeliverReq,     //   "CMPP_DELIVER",
	CMPP_DELIVER_RESP:     NewDeliverResp,    //   "CMPP_DELIVER_RESP",
	CMPP_QUERY:            NewQueryReq,       //   "CMPP_QUERY",
	CMPP_QUERY_RESP:       NewQueryResp,      //   "CMPP_QUERY_RESP",
	CMPP_CANCEL:           NewCancelReq,      //   "CMPP_CANCEL",
	CMPP_CANCEL_RESP:      NewCancelResp,     //   "CMPP_CANCEL_RESP",
	CMPP_ACTIVE_TEST:      NewActiveTestReq,  //   "CMPP_ACTIVE_TEST",
	CMPP_ACTIVE_TEST_RESP: NewActiveTestResp, //  "CMPP_ACTIVE_TEST_RESP",
	CMPP_FWD:              NewFwdReq,         //  "CMPP_FWD",
	CMPP_FWD_RESP:         NewFwdResp,        //   "CMPP_FWD_RESP",
	CMPP_REQUEST_MAX:      nil,               //   "CMPP_REQUEST_MAX",
	CMPP_RESPONSE_MAX:     nil,               //  "CMPP_RESPONSE_MAX",
}

func CreatePDUFromCmdID(cmdID codec.CommandId, ver codec.Version) (codec.PDU, error) {
	if g, ok := pduMap[cmdID]; ok {
		return g(ver), nil
	}
	return nil, smserror.ErrUnknownCommandID
}
