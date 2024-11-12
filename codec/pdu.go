package codec

type CommandId uint32
type Version byte
type CommandStatus uint32

type Header interface {
	AssignSequenceNumber()
	// ResetSequenceNumber resets sequence number.
	ResetSequenceNumber()
	// GetSequenceNumber returns assigned sequence number.
	GetSequenceNumber() int32
	// SetSequenceNumber manually sets sequence number.
	SetSequenceNumber(int32)
}

// PDU represents PDU interface.
type PDU interface {
	// Marshal PDU to buffer.
	Marshal(*BytesWriter)

	// Unmarshal PDU from buffer.
	Unmarshal(*BytesReader) error

	// GetResponse PDU.
	GetResponse() PDU

	// RegisterOptionalParam assigns an optional param.
	RegisterOptionalParam(Field)

	// GetHeader returns PDU header.
	GetHeader() Header

	// IsOk returns true if command status is OK.
	IsOk() bool

	// IsGNack returns true if PDU is GNack.
	IsGNack() bool

	// AssignSequenceNumber assigns sequence number auto-incrementally.
	AssignSequenceNumber()

	// ResetSequenceNumber resets sequence number.
	ResetSequenceNumber()

	// GetSequenceNumber returns assigned sequence number.
	GetSequenceNumber() int32

	// SetSequenceNumber manually sets sequence number.
	SetSequenceNumber(int32)
	String() string
}
