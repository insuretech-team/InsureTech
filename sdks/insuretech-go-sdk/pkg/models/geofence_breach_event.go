package models

import (
	"time"
)

// GeofenceBreachEvent represents a geofence_breach_event
type GeofenceBreachEvent struct {
	GeofenceName string `json:"geofence_name,omitempty"`
	Latitude float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	DeviceId string `json:"device_id,omitempty"`
	DeviceSerial string `json:"device_serial,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
}
