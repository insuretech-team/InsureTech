package models

import (
	"time"
)

// VehicleCreatedEvent represents a vehicle_created_event
type VehicleCreatedEvent struct {
	CreatedAt time.Time `json:"created_at,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`
	Model string `json:"model,omitempty"`
	Price string `json:"price,omitempty"`
	Type *VehicleType `json:"type,omitempty"`
	VehicleId string `json:"vehicle_id,omitempty"`
}
