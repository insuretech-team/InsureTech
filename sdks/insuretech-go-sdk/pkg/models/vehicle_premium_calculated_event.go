package models

import (
	"time"
)

// VehiclePremiumCalculatedEvent represents a vehicle_premium_calculated_event
type VehiclePremiumCalculatedEvent struct {
	BasePremiumAmount string `json:"base_premium_amount,omitempty"`
	CalculatedAt time.Time `json:"calculated_at,omitempty"`
	CalculationId string `json:"calculation_id,omitempty"`
	CompPremium1YearAmount string `json:"comp_premium_1_year_amount,omitempty"`
	EventId string `json:"event_id,omitempty"`
	RegistrationId string `json:"registration_id,omitempty"`
	RegistrationNumber string `json:"registration_number,omitempty"`
	TpPremium1YearAmount string `json:"tp_premium_1_year_amount,omitempty"`
}
