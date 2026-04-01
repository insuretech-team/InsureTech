package models

import (
	"time"
)

// VehicleDeletedEvent represents a vehicle_deleted_event
type VehicleDeletedEvent struct {
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Model string `json:"model,omitempty"`
	Permanent bool `json:"permanent,omitempty"`
	VehicleId string `json:"vehicle_id,omitempty"`
}
