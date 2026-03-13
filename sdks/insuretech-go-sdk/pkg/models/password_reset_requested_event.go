package models

import (
	"time"
)

// PasswordResetRequestedEvent represents a password_reset_requested_event
type PasswordResetRequestedEvent struct {
	EventId string `json:"event_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	DeviceType string `json:"device_type,omitempty"`
}
