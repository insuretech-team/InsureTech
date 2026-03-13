package models

import (
	"time"
)

// EmailLoginFailedEvent represents a email_login_failed_event
type EmailLoginFailedEvent struct {
	FailedAttemptsCount int `json:"failed_attempts_count,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	EventId string `json:"event_id,omitempty"`
	EmailMasked string `json:"email_masked,omitempty"`
	FailureReason string `json:"failure_reason,omitempty"`
}
