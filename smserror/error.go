package smserror

import (
	"errors"
)

type SmsError error

// Errors for conn operations
var (

	// Common errors.
	ErrMethodParamsInvalid SmsError = errors.New("params passed to method is invalid")

	// Protocol errors.
	ErrTotalLengthInvalid    SmsError = errors.New("total_length in Packet data is invalid")
	ErrCommandIdInvalid      SmsError = errors.New("command_Id in Packet data is invalid")
	ErrCommandIdNotSupported SmsError = errors.New("command_Id in Packet data is not supported")

	ErrConnIsClosed       SmsError = errors.New("connection is closed")
	ErrReadCmdIDTimeout   SmsError = errors.New("read commandId timeout")
	ErrReadPktBodyTimeout SmsError = errors.New("read packet body timeout")
	ErrNotCompleted       SmsError = errors.New("data not being handled completed")
	ErrRespNotMatch       SmsError = errors.New("the response is not matched with the request")
	ErrEmptyServerAddr    SmsError = errors.New("sms server listen: empty server addr")
	ErrNoHandlers         SmsError = errors.New("sms server: no connection handler")
	ErrUnsupportedPkt     SmsError = errors.New("sms server read packet: receive a unsupported pkt")
	ErrProtoNotSupport    SmsError = errors.New("sms unsupported proto")
	ErrPktIsNil           SmsError = errors.New("sms packet is nil")
	ErrVersionNotMatch    SmsError = errors.New("sms version not match")
)
