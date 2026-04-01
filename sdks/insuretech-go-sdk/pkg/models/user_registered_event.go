package models

import (
	"time"
)

// UserRegisteredEvent represents a user_registered_event
type UserRegisteredEvent struct {
	DeviceType string `json:"device_type,omitempty"`
	Email string `json:"email,omitempty"`
	EventId string `json:"event_id,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
	Portal string `json:"portal,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
