package models

import (
	"time"
)

// VehicleUpdatedEvent represents a vehicle_updated_event
type VehicleUpdatedEvent struct {
	ChangedFields []string `json:"changed_fields,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Model string `json:"model,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	VehicleId string `json:"vehicle_id,omitempty"`
}
