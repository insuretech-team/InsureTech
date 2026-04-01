package models

import (
	"time"
)

// PricingRuleCreatedEvent represents a pricing_rule_created_event
type PricingRuleCreatedEvent struct {
	CreatedBy string `json:"created_by,omitempty"`
	EffectiveFrom time.Time `json:"effective_from,omitempty"`
	EffectiveTo time.Time `json:"effective_to,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PricingConfigId string `json:"pricing_config_id,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	RuleType string `json:"rule_type,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
