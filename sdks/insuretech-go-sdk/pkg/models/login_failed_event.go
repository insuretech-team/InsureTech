package models

import (
	"time"
)

// LoginFailedEvent represents a login_failed_event
type LoginFailedEvent struct {
	DeviceType string `json:"device_type,omitempty"`
	EventId string `json:"event_id,omitempty"`
	FailedAttemptsCount int `json:"failed_attempts_count,omitempty"`
	FailureReason string `json:"failure_reason,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
