package models


// FraudRuleDeactivationRequest represents a fraud_rule_deactivation_request
type FraudRuleDeactivationRequest struct {
	Reason string `json:"reason,omitempty"`
	RuleId string `json:"rule_id"`
}
