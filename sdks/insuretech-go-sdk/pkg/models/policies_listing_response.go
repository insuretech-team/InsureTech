package models


// PoliciesListingResponse represents a policies_listing_response
type PoliciesListingResponse struct {
	Policies []*Policy `json:"policies,omitempty"`
	Total int `json:"total,omitempty"`
}
