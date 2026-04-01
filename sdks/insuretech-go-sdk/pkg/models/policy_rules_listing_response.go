package models


// PolicyRulesListingResponse represents a policy_rules_listing_response
type PolicyRulesListingResponse struct {
	NextPageToken string `json:"next_page_token,omitempty"`
	Policies []*PolicyRule `json:"policies,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
