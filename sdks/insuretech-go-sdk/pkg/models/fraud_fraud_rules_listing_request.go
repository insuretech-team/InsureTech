package models


// FraudFraudRulesListingRequest represents a fraud_fraud_rules_listing_request
type FraudFraudRulesListingRequest struct {
	ActiveOnly bool `json:"active_only,omitempty"`
	Category string `json:"category"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}
