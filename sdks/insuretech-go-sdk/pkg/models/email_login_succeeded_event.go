package models

import (
	"time"
)

// EmailLoginSucceededEvent represents a email_login_succeeded_event
type EmailLoginSucceededEvent struct {
	UserId string `json:"user_id,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	EmailMasked string `json:"email_masked,omitempty"`
	UserType string `json:"user_type,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	EventId string `json:"event_id,omitempty"`
	DeviceName string `json:"device_name,omitempty"`
}
