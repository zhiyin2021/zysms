package enum

import (
	"errors"
)

type SmsError error

// Errors for conn operations
var (
	ErrConnIsClosed       SmsError = errors.New("connection is closed")
	ErrReadCmdIDTimeout   SmsError = errors.New("read commandId timeout")
	ErrReadPktBodyTimeout SmsError = errors.New("read packet body timeout")
	ErrNotCompleted       SmsError = errors.New("data not being handled completed")
	ErrRespNotMatch       SmsError = errors.New("the response is not matched with the request")
	ErrEmptyServerAddr    SmsError = errors.New("cmpp server listen: empty server addr")
	ErrNoHandlers         SmsError = errors.New("cmpp server: no connection handler")
	ErrUnsupportedPkt     SmsError = errors.New("cmpp server read packet: receive a unsupported pkt")
)
