package models


// PricingRule represents a pricing_rule
type PricingRule struct {
	Action *RuleAction `json:"action,omitempty"`
	Conditions []*RuleCondition `json:"conditions,omitempty"`
	RuleId string `json:"rule_id,omitempty"`
	RuleName string `json:"rule_name,omitempty"`
	Type *RuleType `json:"type,omitempty"`
}
