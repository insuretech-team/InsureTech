package models


// FraudFraudRulesListingRequest represents a fraud_fraud_rules_listing_request
type FraudFraudRulesListingRequest struct {
	Category string `json:"category"`
	ActiveOnly bool `json:"active_only,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}
