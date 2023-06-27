package event

type SmsEvent byte

const (
	SmsEventUnknown SmsEvent = iota

	SmsEventAuthReq
	SmsEventAuthRsp

	SmsEventActiveTestReq
	SmsEventActiveTestRsp

	SmsEventSubmitReq
	SmsEventSubmitRsp

	SmsEventDeliverReq
	SmsEventDeliverRsp

	SmsEventQueryReq
	SmsEventQueryRsp

	SmsEventCancelReq
	SmsEventCancelRsp

	SmsEventReportReq
	SmsEventReportRsp

	SmsEventFwdReq
	SmsEventFwdRsp

	SmsEventTerminateReq
	SmsEventTerminateRsp
)
