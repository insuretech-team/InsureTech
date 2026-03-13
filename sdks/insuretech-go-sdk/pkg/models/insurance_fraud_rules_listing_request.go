package models


// InsuranceFraudRulesListingRequest represents a insurance_fraud_rules_listing_request
type InsuranceFraudRulesListingRequest struct {
	Page int `json:"page"`
	PageSize int `json:"page_size,omitempty"`
}
