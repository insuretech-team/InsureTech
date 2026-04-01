package models

import (
	"time"
)

// FirmwareUpdatedEvent represents a firmware_updated_event
type FirmwareUpdatedEvent struct {
	DeviceId string `json:"device_id,omitempty"`
	DeviceSerial string `json:"device_serial,omitempty"`
	EventId string `json:"event_id,omitempty"`
	NewVersion string `json:"new_version,omitempty"`
	OldVersion string `json:"old_version,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UpdateSuccessful bool `json:"update_successful,omitempty"`
}
