package models

import (
	"time"
)

// QuotingPremiumCalculatedEvent represents a quoting_premium_calculated_event
type QuotingPremiumCalculatedEvent struct {
	CalculatedAt time.Time `json:"calculated_at,omitempty"`
	CalculationId string `json:"calculation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ExecutionTimeMs int `json:"execution_time_ms,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	RuleCount int `json:"rule_count,omitempty"`
	TotalPremium *Money `json:"total_premium,omitempty"`
}
