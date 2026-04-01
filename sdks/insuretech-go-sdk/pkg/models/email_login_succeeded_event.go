package models

import (
	"time"
)

// EmailLoginSucceededEvent represents a email_login_succeeded_event
type EmailLoginSucceededEvent struct {
	DeviceName string `json:"device_name,omitempty"`
	EmailMasked string `json:"email_masked,omitempty"`
	EventId string `json:"event_id,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	UserId string `json:"user_id,omitempty"`
	UserType string `json:"user_type,omitempty"`
}
