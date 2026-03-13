package models


// PolicyRulesListingRequest represents a policy_rules_listing_request
type PolicyRulesListingRequest struct {
	Domain string `json:"domain"`
	Subject string `json:"subject,omitempty"`
	Object string `json:"object,omitempty"`
	ActiveOnly bool `json:"active_only,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}
