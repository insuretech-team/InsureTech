package models


// ClaimsListingRequest represents a claims_listing_request
type ClaimsListingRequest struct {
	CustomerId string `json:"customer_id"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PolicyId string `json:"policy_id"`
}
