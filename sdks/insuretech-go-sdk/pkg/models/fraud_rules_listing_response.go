package models


// FraudRulesListingResponse represents a fraud_rules_listing_response
type FraudRulesListingResponse struct {
	FraudRules []*FraudRule `json:"fraud_rules,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
