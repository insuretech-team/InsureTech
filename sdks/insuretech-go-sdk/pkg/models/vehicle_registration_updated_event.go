package models

import (
	"time"
)

// VehicleRegistrationUpdatedEvent represents a vehicle_registration_updated_event
type VehicleRegistrationUpdatedEvent struct {
	EventId string `json:"event_id,omitempty"`
	NewStatus *RegistrationStatus `json:"new_status,omitempty"`
	OldStatus *RegistrationStatus `json:"old_status,omitempty"`
	RegistrationId string `json:"registration_id,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
