package models


// RuleCondition represents a rule_condition
type RuleCondition struct {
	Field string `json:"field,omitempty"`
	Operator string `json:"operator,omitempty"`
	Value string `json:"value,omitempty"`
}
