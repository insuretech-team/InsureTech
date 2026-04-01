package models

import (
	"time"
)

// SMSDeliveryReportEvent represents a sms_delivery_report_event
type SMSDeliveryReportEvent struct {
	Carrier string `json:"carrier,omitempty"`
	DeliveredAt time.Time `json:"delivered_at,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Msisdn string `json:"msisdn,omitempty"`
	OtpId string `json:"otp_id,omitempty"`
	ProviderMessageId string `json:"provider_message_id,omitempty"`
	Status string `json:"status,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
