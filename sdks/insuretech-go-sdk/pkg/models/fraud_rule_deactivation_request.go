package models


// FraudRuleDeactivationRequest represents a fraud_rule_deactivation_request
type FraudRuleDeactivationRequest struct {
	RuleId string `json:"rule_id"`
	Reason string `json:"reason,omitempty"`
}
