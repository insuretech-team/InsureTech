package models


// UserClaimsListingResponse represents a user_claims_listing_response
type UserClaimsListingResponse struct {
	Claims []*Claim `json:"claims,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
