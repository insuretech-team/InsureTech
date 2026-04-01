package models


// PolicyRulesListingRequest represents a policy_rules_listing_request
type PolicyRulesListingRequest struct {
	ActiveOnly bool `json:"active_only,omitempty"`
	Domain string `json:"domain"`
	Object string `json:"object,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	Subject string `json:"subject,omitempty"`
}
