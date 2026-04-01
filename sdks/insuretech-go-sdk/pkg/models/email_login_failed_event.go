package models

import (
	"time"
)

// EmailLoginFailedEvent represents a email_login_failed_event
type EmailLoginFailedEvent struct {
	EmailMasked string `json:"email_masked,omitempty"`
	EventId string `json:"event_id,omitempty"`
	FailedAttemptsCount int `json:"failed_attempts_count,omitempty"`
	FailureReason string `json:"failure_reason,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}
