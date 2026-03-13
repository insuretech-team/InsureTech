package models


// ClaimsListingResponse represents a claims_listing_response
type ClaimsListingResponse struct {
	Claims []*Claim `json:"claims,omitempty"`
	Total int `json:"total,omitempty"`
}
