package models

import (
	"time"
)

// PremiumCalculatedEvent represents a premium_calculated_event
type PremiumCalculatedEvent struct {
	ProductId string `json:"product_id,omitempty"`
	BasePremium *Money `json:"base_premium,omitempty"`
	FinalPremium *Money `json:"final_premium,omitempty"`
	AppliedRules []string `json:"applied_rules,omitempty"`
	InputFactors map[string]interface{} `json:"input_factors,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	CalculatedForUser string `json:"calculated_for_user,omitempty"`
	EventId string `json:"event_id,omitempty"`
}
