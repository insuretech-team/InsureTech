package models


// PoliciesListingRequest represents a policies_listing_request
type PoliciesListingRequest struct {
	CustomerId string `json:"customer_id"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	TenantId string `json:"tenant_id"`
}
