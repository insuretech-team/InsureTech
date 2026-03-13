package models


// PolicyRulesListingResponse represents a policy_rules_listing_response
type PolicyRulesListingResponse struct {
	Error *Error `json:"error,omitempty"`
	Policies []*PolicyRule `json:"policies,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
