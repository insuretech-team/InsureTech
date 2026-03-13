package models


// FraudRuleUpdateResponse represents a fraud_rule_update_response
type FraudRuleUpdateResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
