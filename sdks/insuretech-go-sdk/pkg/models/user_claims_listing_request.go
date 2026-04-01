package models


// UserClaimsListingRequest represents a user_claims_listing_request
type UserClaimsListingRequest struct {
	CustomerId string `json:"customer_id"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	Status *ClaimStatus `json:"status,omitempty"`
}
