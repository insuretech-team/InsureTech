package models


// FraudFraudRuleUpdateRequest represents a fraud_fraud_rule_update_request
type FraudFraudRuleUpdateRequest struct {
	FraudRule *FraudRule `json:"fraud_rule,omitempty"`
	RuleId string `json:"rule_id"`
}
