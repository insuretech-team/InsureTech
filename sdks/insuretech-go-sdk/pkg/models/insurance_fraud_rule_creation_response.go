package models


// InsuranceFraudRuleCreationResponse represents a insurance_fraud_rule_creation_response
type InsuranceFraudRuleCreationResponse struct {
	Rule *FraudRule `json:"rule,omitempty"`
}
