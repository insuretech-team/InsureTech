package models

import (
	"time"
)

// VehiclePremiumCalculation represents a vehicle_premium_calculation
type VehiclePremiumCalculation struct {
	AccidentalCover bool `json:"accidental_cover,omitempty"`
	AgeMultiplier float64 `json:"age_multiplier,omitempty"`
	BasePremium *Money `json:"base_premium,omitempty"`
	CalculatedAt time.Time `json:"calculated_at"`
	CalculationDurationMs int `json:"calculation_duration_ms,omitempty"`
	CalculationId string `json:"calculation_id"`
	CompPremium1Year *Money `json:"comp_premium_1_year,omitempty"`
	CompPremium2Year *Money `json:"comp_premium_2_year,omitempty"`
	CompPremium3Year *Money `json:"comp_premium_3_year,omitempty"`
	LocationMultiplier float64 `json:"location_multiplier,omitempty"`
	RegistrationId string `json:"registration_id"`
	TpPremium1Year *Money `json:"tp_premium_1_year,omitempty"`
	TpPremium2Year *Money `json:"tp_premium_2_year,omitempty"`
	TpPremium3Year *Money `json:"tp_premium_3_year,omitempty"`
	TypeMultiplier float64 `json:"type_multiplier,omitempty"`
	ValueMultiplier float64 `json:"value_multiplier,omitempty"`
}
