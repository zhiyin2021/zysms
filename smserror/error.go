package smserror

import (
	"errors"
	"fmt"
)

type SmsError struct {
	Err  error
	Code int
}

func (e SmsError) Error() string {
	return fmt.Sprintf("code:%d,msg:%s", e.Code, e.Err.Error())
}
func NewSmsErr(code int, msg string) *SmsError {
	return &SmsError{Code: code, Err: errors.New(msg)}
}

// Errors for conn operations
var (

	// Common errors.
	ErrMethodParamsInvalid = NewSmsErr(1, "params passed to method is invalid")

	// Protocol errors.
	ErrTotalLengthInvalid    = NewSmsErr(2, "total_length in Packet data is invalid")
	ErrCommandIdInvalid      = NewSmsErr(3, "command_Id in Packet data is invalid")
	ErrCommandIdNotSupported = NewSmsErr(4, "command_Id in Packet data is not supported")

	ErrConnIsClosed       = NewSmsErr(5, "connection is closed")
	ErrReadCmdIDTimeout   = NewSmsErr(6, "read commandId timeout")
	ErrReadPktBodyTimeout = NewSmsErr(7, "read packet body timeout")
	ErrNotCompleted       = NewSmsErr(8, "data not being handled completed")
	ErrRespNotMatch       = NewSmsErr(9, "the response is not matched with the request")
	ErrEmptyServerAddr    = NewSmsErr(10, "sms server listen: empty server addr")
	ErrNoHandlers         = NewSmsErr(11, "sms server: no connection handler")
	ErrUnsupportedPkt     = NewSmsErr(12, "sms server read packet: receive a unsupported pkt")
	ErrProtoNotSupport    = NewSmsErr(13, "sms unsupported proto")
	ErrPktIsNil           = NewSmsErr(14, "sms packet is nil")
	ErrVersionNotMatch    = NewSmsErr(15, "sms version not match")

	// ErrInvalidPDU indicates invalid pdu payload.
	ErrInvalidPDU = NewSmsErr(16, "PDU payload is invalid")

	// ErrUnknownCommandID indicates unknown command id.
	ErrUnknownCommandID = NewSmsErr(17, "unknown command id")

	// ErrWrongDateFormat indicates wrong date format.
	ErrWrongDateFormat = NewSmsErr(18, "wrong date format")

	// ErrShortMessageLengthTooLarge indicates short message length is too large.
	ErrShortMessageLengthTooLarge = NewSmsErr(19, "encoded short message data exceeds size out of range")

	// ErrUDHTooLong UDH-L is larger than total length of short message data
	ErrUDHTooLong = NewSmsErr(20, "user Data Header is too long for PDU short message")
	// Errors for connect resp status.

	ErrnoConnInvalidStruct  uint8 = 1
	ErrnoConnInvalidSrcAddr uint8 = 2
	ErrnoConnAuthFailed     uint8 = 3
	ErrnoConnVerTooHigh     uint8 = 4
	ErrnoConnOthers         uint8 = 5

	ConnRspStatusErrMap = map[uint8]error{
		ErrnoConnInvalidStruct:  errConnInvalidStruct,
		ErrnoConnInvalidSrcAddr: errConnInvalidSrcAddr,
		ErrnoConnAuthFailed:     errConnAuthFailed,
		ErrnoConnVerTooHigh:     errConnVerTooHigh,
		ErrnoConnOthers:         errConnOthers,
	}

	errConnInvalidStruct  = NewSmsErr(1, "connect response status: invalid protocol structure")
	errConnInvalidSrcAddr = NewSmsErr(2, "connect response status: invalid source address")
	errConnAuthFailed     = NewSmsErr(3, "connect response status: auth failed")
	errConnVerTooHigh     = NewSmsErr(4, "connect response status: protocol version is too high")
	errConnOthers         = NewSmsErr(5, "connect response status: other errors")
)
