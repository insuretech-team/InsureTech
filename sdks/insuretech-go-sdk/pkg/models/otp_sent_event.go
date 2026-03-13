package models

import (
	"time"
)

// OTPSentEvent represents a otp_sent_event
type OTPSentEvent struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	SenderId string `json:"sender_id,omitempty"`
	MaskingUsed bool `json:"masking_used,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Channel string `json:"channel,omitempty"`
	Provider string `json:"provider,omitempty"`
	ProviderMessageId string `json:"provider_message_id,omitempty"`
	OtpId string `json:"otp_id,omitempty"`
	Type string `json:"type,omitempty"`
}
