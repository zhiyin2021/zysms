package enum

// Conn States
const (
	CONN_DISCONNECTED int32 = iota
	CONN_CONNECTED
	CONN_AUTHOK
)

// Stat
const (
	REPORT_DELIVERED = "DELIVRD"
	REPORT_EXPIRED   = "EXPIRED"
	REPORT_DELETED   = "DELETED"
	REPORT_UNDELIV   = "UNDELIV"
	REPORT_ACCEPTD   = "ACCEPTD"
	REPORT_UNKNOWN   = "UNKNOWN"
	REPORT_REJECTD   = "REJECTD"
	REPORT_MA        = "MA:"
	REPORT_MB        = "MB:"
	REPORT_CA        = "CA:"
	REPORT_CB        = "CB:"
)
