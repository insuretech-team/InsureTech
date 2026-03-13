package models

import (
	"time"
)

// FirmwareUpdatedEvent represents a firmware_updated_event
type FirmwareUpdatedEvent struct {
	DeviceId string `json:"device_id,omitempty"`
	DeviceSerial string `json:"device_serial,omitempty"`
	OldVersion string `json:"old_version,omitempty"`
	NewVersion string `json:"new_version,omitempty"`
	UpdateSuccessful bool `json:"update_successful,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
}
