package models


// InsuranceFraudRulesListingResponse represents a insurance_fraud_rules_listing_response
type InsuranceFraudRulesListingResponse struct {
	Rules []*FraudRule `json:"rules,omitempty"`
	Total int `json:"total,omitempty"`
}
