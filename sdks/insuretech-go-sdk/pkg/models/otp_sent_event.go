package models

import (
	"time"
)

// OTPSentEvent represents a otp_sent_event
type OTPSentEvent struct {
	Channel string `json:"channel,omitempty"`
	EventId string `json:"event_id,omitempty"`
	MaskingUsed bool `json:"masking_used,omitempty"`
	OtpId string `json:"otp_id,omitempty"`
	Provider string `json:"provider,omitempty"`
	ProviderMessageId string `json:"provider_message_id,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	SenderId string `json:"sender_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Type string `json:"type,omitempty"`
}
