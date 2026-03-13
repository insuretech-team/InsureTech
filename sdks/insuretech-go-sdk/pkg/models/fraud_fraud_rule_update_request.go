package models


// FraudFraudRuleUpdateRequest represents a fraud_fraud_rule_update_request
type FraudFraudRuleUpdateRequest struct {
	RuleId string `json:"rule_id"`
	FraudRule *FraudRule `json:"fraud_rule,omitempty"`
}
