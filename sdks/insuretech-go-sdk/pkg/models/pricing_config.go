package models

import (
	"time"
)

// PricingConfig represents a pricing_config
type PricingConfig struct {
	CreatedAt time.Time `json:"created_at,omitempty"`
	EffectiveFrom time.Time `json:"effective_from,omitempty"`
	EffectiveTo time.Time `json:"effective_to,omitempty"`
	PricingConfigId string `json:"pricing_config_id,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	Rules []*ProductsPricingRule `json:"rules,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
