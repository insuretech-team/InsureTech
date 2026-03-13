package models


// FraudRuleUpdateRequest represents a fraud_rule_update_request
type FraudRuleUpdateRequest struct {
	RuleId string `json:"rule_id"`
	FraudRule *FraudRule `json:"fraud_rule,omitempty"`
}
