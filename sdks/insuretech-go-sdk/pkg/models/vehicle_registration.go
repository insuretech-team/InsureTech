package models

import (
	"time"
)

// VehicleRegistration represents a vehicle_registration
type VehicleRegistration struct {
	AdditionalInfo map[string]interface{} `json:"additional_info,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	CurrentValue string `json:"current_value,omitempty"`
	OwnerId string `json:"owner_id"`
	RegistrationId string `json:"registration_id"`
	RegistrationNumber string `json:"registration_number"`
	RegistrationState string `json:"registration_state"`
	RegistrationYear int `json:"registration_year"`
	Status interface{} `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	VehicleAge int `json:"vehicle_age,omitempty"`
	VehicleId string `json:"vehicle_id"`
}
