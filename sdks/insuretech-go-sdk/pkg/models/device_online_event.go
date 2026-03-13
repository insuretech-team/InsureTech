package models

import (
	"time"
)

// DeviceOnlineEvent represents a device_online_event
type DeviceOnlineEvent struct {
	DeviceSerial string `json:"device_serial,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	OfflineDurationMinutes int `json:"offline_duration_minutes,omitempty"`
	EventId string `json:"event_id,omitempty"`
	DeviceId string `json:"device_id,omitempty"`
}
