package models


// ProductsPricingRule represents a products_pricing_rule
type ProductsPricingRule struct {
	Action *RuleAction `json:"action,omitempty"`
	Conditions []*RuleCondition `json:"conditions,omitempty"`
	RuleId string `json:"rule_id,omitempty"`
	RuleName string `json:"rule_name,omitempty"`
	Type *RuleType `json:"type,omitempty"`
}
