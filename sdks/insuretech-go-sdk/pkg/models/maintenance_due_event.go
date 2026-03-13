package models

import (
	"time"
)

// MaintenanceDueEvent represents a maintenance_due_event
type MaintenanceDueEvent struct {
	DeviceId string `json:"device_id,omitempty"`
	DeviceSerial string `json:"device_serial,omitempty"`
	MaintenanceType string `json:"maintenance_type,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
}
