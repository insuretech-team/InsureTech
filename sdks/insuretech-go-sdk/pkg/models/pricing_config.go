package models

import (
	"time"
)

// PricingConfig represents a pricing_config
type PricingConfig struct {
	ProductId string `json:"product_id,omitempty"`
	Rules []*PricingRule `json:"rules,omitempty"`
	EffectiveFrom time.Time `json:"effective_from,omitempty"`
	EffectiveTo time.Time `json:"effective_to,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	PricingConfigId string `json:"pricing_config_id,omitempty"`
}
