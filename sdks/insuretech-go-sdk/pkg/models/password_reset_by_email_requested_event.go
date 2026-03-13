package models

import (
	"time"
)

// PasswordResetByEmailRequestedEvent represents a password_reset_by_email_requested_event
type PasswordResetByEmailRequestedEvent struct {
	OtpId string `json:"otp_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	EventId string `json:"event_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	EmailMasked string `json:"email_masked,omitempty"`
}
