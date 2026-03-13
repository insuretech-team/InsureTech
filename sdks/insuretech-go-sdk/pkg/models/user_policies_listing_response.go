package models


// UserPoliciesListingResponse represents a user_policies_listing_response
type UserPoliciesListingResponse struct {
	Policies []*Policy `json:"policies,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
