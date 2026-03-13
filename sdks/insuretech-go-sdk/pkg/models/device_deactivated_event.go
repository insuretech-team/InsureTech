package models

import (
	"time"
)

// DeviceDeactivatedEvent represents a device_deactivated_event
type DeviceDeactivatedEvent struct {
	EventId string `json:"event_id,omitempty"`
	DeviceId string `json:"device_id,omitempty"`
	DeviceSerial string `json:"device_serial,omitempty"`
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
