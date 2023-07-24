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

func CreatePDUHeader(header Header, ver codec.Version) (codec.PDU, error) {
	base := newBase(ver, header.CommandID, header.SequenceNumber)
	switch header.CommandID {
	case CMPP_CONNECT:
		return &ConnReq{base: base}, nil
	case CMPP_CONNECT_RESP:
		return &ConnResp{base: base}, nil
	case CMPP_TERMINATE:
		return &TerminateReq{base: base}, nil
	case CMPP_TERMINATE_RESP:
		return &TerminateResp{base: base}, nil
	case CMPP_SUBMIT:
		return &SubmitReq{base: base}, nil
	case CMPP_SUBMIT_RESP:
		return &SubmitResp{base: base}, nil
	case CMPP_DELIVER:
		return &DeliverReq{base: base}, nil
	case CMPP_DELIVER_RESP:
		return &DeliverResp{base: base}, nil
	case CMPP_QUERY:
		return &QueryReq{base: base}, nil
	case CMPP_QUERY_RESP:
		return &QueryResp{base: base}, nil
	case CMPP_CANCEL:
		return &CancelReq{base: base}, nil
	case CMPP_CANCEL_RESP:
		return &CancelResp{base: base}, nil
	case CMPP_ACTIVE_TEST:
		return &ActiveTestReq{base: base}, nil
	case CMPP_ACTIVE_TEST_RESP:
		return &ActiveTestResp{base: base}, nil
	case CMPP_FWD:
		return &FwdReq{base: base}, nil
	case CMPP_FWD_RESP:
		return &FwdResp{base: base}, nil
	default:
		return nil, smserror.ErrUnknownCommandID
	}
}
