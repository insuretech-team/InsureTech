package models

import (
	"time"
)

// EmailVerificationSentEvent represents a email_verification_sent_event
type EmailVerificationSentEvent struct {
	UserId string `json:"user_id,omitempty"`
	EmailMasked string `json:"email_masked,omitempty"`
	OtpId string `json:"otp_id,omitempty"`
	Type string `json:"type,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	EventId string `json:"event_id,omitempty"`
}
