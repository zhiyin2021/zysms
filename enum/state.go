package enum

type State uint8

// Conn States
const (
	CONN_CLOSED State = iota
	CONN_CONNECTED
	CONN_AUTHOK
)
