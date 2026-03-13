package models


// FraudRuleCreationResponse represents a fraud_rule_creation_response
type FraudRuleCreationResponse struct {
	Error *Error `json:"error,omitempty"`
	RuleId string `json:"rule_id,omitempty"`
	Message string `json:"message,omitempty"`
}
