package models

import (
	"time"
)

// DeviceRegisteredEvent represents a device_registered_event
type DeviceRegisteredEvent struct {
	DeviceId string `json:"device_id,omitempty"`
	DeviceSerial string `json:"device_serial,omitempty"`
	DeviceType string `json:"device_type,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`
	Model string `json:"model,omitempty"`
	OwnerId string `json:"owner_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
