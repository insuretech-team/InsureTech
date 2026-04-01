package models

import (
	"time"
)

// PasswordResetByEmailRequestedEvent represents a password_reset_by_email_requested_event
type PasswordResetByEmailRequestedEvent struct {
	EmailMasked string `json:"email_masked,omitempty"`
	EventId string `json:"event_id,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	OtpId string `json:"otp_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
