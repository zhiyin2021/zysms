package smgp

import (
	"fmt"
)

// CmppErr indicates smpp error(s), compatible with OpenSMPP.
type SmgpErr struct {
	err              string
	serialVersionUID int64
}

// Error interface.
func (s *SmgpErr) Error() string {
	return fmt.Sprintf("Error happened: [%s]. SerialVersionUID: [%d]", s.err, s.serialVersionUID)
}

var ()
