package models


// FraudRuleCreationRequest represents a fraud_rule_creation_request
type FraudRuleCreationRequest struct {
	FraudRule *FraudRule `json:"fraud_rule"`
}
