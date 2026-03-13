package models


// FraudFraudRuleCreationResponse represents a fraud_fraud_rule_creation_response
type FraudFraudRuleCreationResponse struct {
	RuleId string `json:"rule_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
