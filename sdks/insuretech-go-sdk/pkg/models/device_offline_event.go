package models

import (
	"time"
)

// DeviceOfflineEvent represents a device_offline_event
type DeviceOfflineEvent struct {
	DeviceId string `json:"device_id,omitempty"`
	DeviceSerial string `json:"device_serial,omitempty"`
	EventId string `json:"event_id,omitempty"`
	LastSeenAt time.Time `json:"last_seen_at,omitempty"`
	OfflineDurationMinutes int `json:"offline_duration_minutes,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
