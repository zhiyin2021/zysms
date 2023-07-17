package cmpp

import "fmt"

// CmppErr indicates smpp error(s), compatible with OpenSMPP.
type CmppErr struct {
	err              string
	serialVersionUID int64
}

// Error interface.
func (s *CmppErr) Error() string {
	return fmt.Sprintf("Error happened: [%s]. SerialVersionUID: [%d]", s.err, s.serialVersionUID)
}

var (
	// ErrInvalidPDU indicates invalid pdu payload.
	ErrInvalidPDU error = &CmppErr{err: "PDU payload is invalid", serialVersionUID: -6985061862208729984}

	// ErrUnknownCommandID indicates unknown command id.
	ErrUnknownCommandID error = &CmppErr{err: "Unknown command id", serialVersionUID: -5091873576710864441}

	// ErrWrongDateFormat indicates wrong date format.
	ErrWrongDateFormat error = &CmppErr{err: "Wrong date format", serialVersionUID: 5831937612139037591}

	// ErrShortMessageLengthTooLarge indicates short message length is too large.
	ErrShortMessageLengthTooLarge error = &CmppErr{err: fmt.Sprintf("Encoded short message data exceeds size of %d", SM_MSG_LEN), serialVersionUID: 78237205927624}

	// ErrUDHTooLong UDH-L is larger than total length of short message data
	ErrUDHTooLong = fmt.Errorf("User Data Header is too long for PDU short message")
)
