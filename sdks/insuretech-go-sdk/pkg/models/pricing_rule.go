package models


// PricingRule represents a pricing_rule
type PricingRule struct {
	Type *RuleType `json:"type,omitempty"`
	Conditions []*RuleCondition `json:"conditions,omitempty"`
	Action *RuleAction `json:"action,omitempty"`
	RuleId string `json:"rule_id,omitempty"`
	RuleName string `json:"rule_name,omitempty"`
}
