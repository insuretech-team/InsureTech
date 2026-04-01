package models

import (
	"time"
)

// PremiumCalculatedEvent represents a premium_calculated_event
type PremiumCalculatedEvent struct {
	AppliedRules []string `json:"applied_rules,omitempty"`
	BasePremium *Money `json:"base_premium,omitempty"`
	CalculatedForUser string `json:"calculated_for_user,omitempty"`
	EventId string `json:"event_id,omitempty"`
	FinalPremium *Money `json:"final_premium,omitempty"`
	InputFactors map[string]interface{} `json:"input_factors,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
