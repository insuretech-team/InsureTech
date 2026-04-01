package models

import (
	"time"
)

// UserLoggedOutEvent represents a user_logged_out_event
type UserLoggedOutEvent struct {
	DeviceType string `json:"device_type,omitempty"`
	EventId string `json:"event_id,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	LogoutReason string `json:"logout_reason,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	SessionType string `json:"session_type,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
