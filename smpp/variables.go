package smpp

// CommandId 命令定义
type CommandId uint32

const (
	_, GENERIC_NACK CommandId = iota, 0x80000000 + iota
	BIND_RECEIVER, BIND_RECEIVER_RESP
	BIND_TRANSMITTER, BIND_TRANSMITTER_RESP
	QUERY_SM, QUERY_SM_RESP
	SUBMIT_SM, SUBMIT_SM_RESP
	DELIVER_SM, DELIVER_SM_RESP
	UNBIND, UNBIND_RESP
	REPLACE_SM, REPLACE_SM_RESP
	CANCEL_SM, CANCEL_SM_RESP
	BIND_TRANSCEIVER, BIND_TRANSCEIVER_RESP
)

type CommandStatus uint32

const (
	ESME_ROK           CommandStatus = iota // No Error
	ESME_RINVMSGLEN                         // Message Length is invalid
	ESME_RINVCMDLEN                         // Command Length is invalid
	ESME_RINVCMDID                          // Invalid Command ID
	ESME_RINVBNDSTS                         // Incorrect BIND Status for given com- mand
	ESME_RALYBND                            // ESME Already in Bound State
	ESME_RINVPRTFLG                         // Invalid Priority Flag
	ESME_RINVREGDLVFLG                      // Invalid Registered Delivery Flag
	ESME_RSYSERR                            // System Error
	_                                       // Reserved
	ESME_RINVSRCADR                         // Invalid Source Address
	ESME_RINVDSTADR                         // Invalid Dest Addr
	ESME_RINVMSGID                          // Message ID is invalid
	ESME_RBINDFAIL                          // Bind Failed
	ESME_RINVPASWD                          // Invalid Password
	ESME_RINVSYSID                          // Invalid System ID
	_                                       // Reserved
	ESME_RCANCELFAIL                        // Cancel SM Failed
	_                                       // Reserved
	_                                       // Reserved
	ESME_RREPLACEFAIL                       // Replace SM Failed
)
