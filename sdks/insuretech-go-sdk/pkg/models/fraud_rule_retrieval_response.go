package models


// FraudRuleRetrievalResponse represents a fraud_rule_retrieval_response
type FraudRuleRetrievalResponse struct {
	Rule *FraudRule `json:"rule,omitempty"`
}
