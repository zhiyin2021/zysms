package smpp

import "github.com/zhiyin2021/zysms/codec"

// AlertNotification PDU is sent by the SMSC to the ESME, when the SMSC has detected that
// a particular mobile subscriber has become available and a delivery pending flag had been
// set for that subscriber from a previous data_sm operation.
type AlertNotification struct {
	base
	SourceAddr Address
	EsmeAddr   Address
}

// NewAlertNotification create new alert notification pdu.
func NewAlertNotification() codec.PDU {
	a := &AlertNotification{
		base: newBase(ALERT_NOTIFICATION, 0),
	}
	return a
}

// GetResponse implements PDU interface.
func (a *AlertNotification) GetResponse() codec.PDU {
	return nil
}

// Marshal implements PDU interface.
func (a *AlertNotification) Marshal(b *codec.BytesWriter) {
	a.base.marshal(b, func(b *codec.BytesWriter) {
		a.SourceAddr.Marshal(b)
		a.EsmeAddr.Marshal(b)
	})
}

// Unmarshal implements PDU interface.
func (a *AlertNotification) Unmarshal(b *codec.BytesReader) error {
	return a.base.unmarshal(b, func(b *codec.BytesReader) (err error) {
		if err = a.SourceAddr.Unmarshal(b); err == nil {
			err = a.EsmeAddr.Unmarshal(b)
		}
		return
	})
}
