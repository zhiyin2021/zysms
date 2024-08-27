package codec

// Source code in this file is copied from: https://github.com/fiorix
import (
	"encoding/binary"
	"encoding/hex"
)

// Tag is the tag of a Tag-Length-Value (TLV) field.
type Tag uint16

// Hex returns hexadecimal representation of tag
func (t Tag) Hex() string {
	var bin [2]byte
	binary.BigEndian.PutUint16(bin[:], uint16(t))
	return hex.EncodeToString(bin[:])
}

// Common Tag-Length-Value (TLV) tags.
// 0002   0001   40     #   TP_udhi
// 0009   0001   04     #   pkTotal
// 000a   0001   01     #   pkNumber
const (
	TagTPUdhi                   Tag = 0x0002
	TagDestAddrSubunit          Tag = 0x0005
	TagDestNetworkType          Tag = 0x0006
	TagDestBearerType           Tag = 0x0007
	TagDestTelematicsID         Tag = 0x0008
	TagPkTotal                  Tag = 0x0009
	TagPkNumber                 Tag = 0x000A
	TagSourceAddrSubunit        Tag = 0x000D
	TagSourceNetworkType        Tag = 0x000E
	TagSourceBearerType         Tag = 0x000F
	TagSourceTelematicsID       Tag = 0x0010
	TagQosTimeToLive            Tag = 0x0017
	TagPayloadType              Tag = 0x0019
	TagAdditionalStatusInfoText Tag = 0x001D
	TagReceiptedMessageID       Tag = 0x001E
	TagMsMsgWaitFacilities      Tag = 0x0030
	TagPrivacyIndicator         Tag = 0x0201
	TagSourceSubaddress         Tag = 0x0202
	TagDestSubaddress           Tag = 0x0203
	TagUserMessageReference     Tag = 0x0204
	TagUserResponseCode         Tag = 0x0205
	TagSourcePort               Tag = 0x020A
	TagDestinationPort          Tag = 0x020B
	TagSarMsgRefNum             Tag = 0x020C
	TagLanguageIndicator        Tag = 0x020D
	TagSarTotalSegments         Tag = 0x020E
	TagSarSegmentSeqnum         Tag = 0x020F
	TagCallbackNumPresInd       Tag = 0x0302
	TagCallbackNumAtag          Tag = 0x0303
	TagNumberOfMessages         Tag = 0x0304
	TagCallbackNum              Tag = 0x0381
	TagDpfResult                Tag = 0x0420
	TagSetDpf                   Tag = 0x0421
	TagMsAvailabilityStatus     Tag = 0x0422
	TagNetworkErrorCode         Tag = 0x0423
	TagMessagePayload           Tag = 0x0424
	TagDeliveryFailureReason    Tag = 0x0425
	TagMoreMessagesToSend       Tag = 0x0426
	TagMessageStateOption       Tag = 0x0427
	TagUssdServiceOp            Tag = 0x0501
	TagDisplayTime              Tag = 0x1201
	TagSmsSignal                Tag = 0x1203
	TagMsValidity               Tag = 0x1204
	TagAlertOnMessageDelivery   Tag = 0x130C
	TagItsReplyType             Tag = 0x1380
	TagItsSessionInfo           Tag = 0x1383
	/*
		消息类型
		0TP (for a one-time password)
		MKT (for a marketing message)
		ARN (for an alert, reminder, or notification)
	*/
	TagMessageType Tag = 0x1414
)

// Field is a PDU Tag-Length-Value (TLV) field
type Field struct {
	Tag  Tag
	Data []byte
}

func NewTlv(tag Tag, data []byte) Field {
	return Field{Tag: tag, Data: data}
}

// String implements the Data interface.
func (t *Field) String() string {
	if l := len(t.Data); l > 0 && t.Data[l-1] == 0x00 {
		return string(t.Data[:l-1])
	}
	return string(t.Data)
}

// String implements the Data interface.
func (t *Field) UInt64() uint64 {
	switch len(t.Data) {
	case 1:
		return uint64(t.Data[0])
	case 2:
		return uint64(binary.BigEndian.Uint16(t.Data))
	case 4:
		return uint64(binary.BigEndian.Uint32(t.Data))
	case 8:
		return binary.BigEndian.Uint64(t.Data)
	default:
		return 0
	}
}

// Marshal to writer.
func (t *Field) Marshal(w *BytesWriter) {
	if len(t.Data) > 0 {
		w.Grow(4 + len(t.Data))
		w.WriteU16(uint16(t.Tag))
		w.WriteU16(uint16(len(t.Data)))
		_, _ = w.Write(t.Data)
	}
}

// Unmarshal from reader.
func (t *Field) Unmarshal(b *BytesReader) (err error) {
	t.Tag = Tag(b.ReadU16())
	ln := b.ReadU16()
	t.Data = b.ReadN(int(ln))
	return b.Err()
}
