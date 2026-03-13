package models


// FraudRuleDeactivationResponse represents a fraud_rule_deactivation_response
type FraudRuleDeactivationResponse struct {
	Error *Error `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}
