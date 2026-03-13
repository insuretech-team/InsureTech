package models


// FraudFraudRuleUpdateResponse represents a fraud_fraud_rule_update_response
type FraudFraudRuleUpdateResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
