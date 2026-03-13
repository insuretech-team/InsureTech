package models


// FraudRuleActivationResponse represents a fraud_rule_activation_response
type FraudRuleActivationResponse struct {
	Error *Error `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}
