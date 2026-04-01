package models

import (
	"time"
)

// VehicleRegisteredEvent represents a vehicle_registered_event
type VehicleRegisteredEvent struct {
	EventId string `json:"event_id,omitempty"`
	OwnerId string `json:"owner_id,omitempty"`
	RegisteredAt time.Time `json:"registered_at,omitempty"`
	RegistrationId string `json:"registration_id,omitempty"`
	RegistrationNumber string `json:"registration_number,omitempty"`
	RegistrationState string `json:"registration_state,omitempty"`
	RegistrationYear int `json:"registration_year,omitempty"`
	VehicleId string `json:"vehicle_id,omitempty"`
}
