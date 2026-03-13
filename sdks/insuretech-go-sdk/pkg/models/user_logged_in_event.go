package models

import (
	"time"
)

// UserLoggedInEvent represents a user_logged_in_event
type UserLoggedInEvent struct {
	SessionType string `json:"session_type,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	DeviceType string `json:"device_type,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	EventId string `json:"event_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	SessionId string `json:"session_id,omitempty"`
}
