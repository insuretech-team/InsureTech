package models


// RuleAction represents a rule_action
type RuleAction struct {
	Type *ActionType `json:"type,omitempty"`
	Value float64 `json:"value,omitempty"`
}
