package models

import (
	"time"
)

// PricingRuleUpdatedEvent represents a pricing_rule_updated_event
type PricingRuleUpdatedEvent struct {
	EventId string `json:"event_id,omitempty"`
	PricingConfigId string `json:"pricing_config_id,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
}
