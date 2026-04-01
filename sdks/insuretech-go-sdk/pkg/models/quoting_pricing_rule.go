package models

import (
	"time"
)

// QuotingPricingRule represents a quoting_pricing_rule
type QuotingPricingRule struct {
	CalculationExpression string `json:"calculation_expression"`
	ConditionExpression string `json:"condition_expression,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	IsActive bool `json:"is_active"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	PricingRuleId string `json:"pricing_rule_id"`
	Priority int `json:"priority"`
	ProductId string `json:"product_id"`
	RuleName string `json:"rule_name"`
	RuleType *PricingRuleType `json:"rule_type"`
	UpdatedAt time.Time `json:"updated_at"`
	ValidFrom time.Time `json:"valid_from"`
	ValidUntil time.Time `json:"valid_until,omitempty"`
}
