package models


// UserPoliciesListingRequest represents a user_policies_listing_request
type UserPoliciesListingRequest struct {
	CustomerId string `json:"customer_id"`
	Status *PolicyStatus `json:"status,omitempty"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}
