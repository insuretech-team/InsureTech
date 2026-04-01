package models

import (
	"time"
)

// DeviceOnlineEvent represents a device_online_event
type DeviceOnlineEvent struct {
	DeviceId string `json:"device_id,omitempty"`
	DeviceSerial string `json:"device_serial,omitempty"`
	EventId string `json:"event_id,omitempty"`
	OfflineDurationMinutes int `json:"offline_duration_minutes,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
