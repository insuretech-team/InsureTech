package models

import (
	"time"
)

// LifePremiumCalculatedEvent represents a life_premium_calculated_event
type LifePremiumCalculatedEvent struct {
	AgeAddition string `json:"age_addition,omitempty"`
	AgeAtEntry int `json:"age_at_entry,omitempty"`
	BasePremium string `json:"base_premium,omitempty"`
	CalculatedAt time.Time `json:"calculated_at,omitempty"`
	CalculationDurationMs int `json:"calculation_duration_ms,omitempty"`
	CalculationId string `json:"calculation_id,omitempty"`
	ConditionMultiplier float64 `json:"condition_multiplier,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	TotalPremium string `json:"total_premium,omitempty"`
}
