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

	// ErrInvalidPDU indicates invalid pdu payload.
	ErrInvalidPDU SmsError = errors.New("PDU payload is invalid")

	// ErrUnknownCommandID indicates unknown command id.
	ErrUnknownCommandID SmsError = errors.New("unknown command id")

	// ErrWrongDateFormat indicates wrong date format.
	ErrWrongDateFormat SmsError = errors.New("wrong date format")

	// ErrShortMessageLengthTooLarge indicates short message length is too large.
	ErrShortMessageLengthTooLarge SmsError = errors.New("encoded short message data exceeds size out of range")

	// ErrUDHTooLong UDH-L is larger than total length of short message data
	ErrUDHTooLong SmsError = errors.New("user Data Header is too long for PDU short message")
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

	errConnInvalidStruct  = errors.New("connect response status: invalid protocol structure")
	errConnInvalidSrcAddr = errors.New("connect response status: invalid source address")
	errConnAuthFailed     = errors.New("connect response status: auth failed")
	errConnVerTooHigh     = errors.New("connect response status: protocol version is too high")
	errConnOthers         = errors.New("connect response status: other errors")
)
