package models

import (
	"time"
)

// EmailVerificationSentEvent represents a email_verification_sent_event
type EmailVerificationSentEvent struct {
	EmailMasked string `json:"email_masked,omitempty"`
	EventId string `json:"event_id,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	OtpId string `json:"otp_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Type string `json:"type,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
