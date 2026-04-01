package models

import (
	"time"
)

// UserLoggedInEvent represents a user_logged_in_event
type UserLoggedInEvent struct {
	DeviceType string `json:"device_type,omitempty"`
	EventId string `json:"event_id,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	SessionType string `json:"session_type,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
