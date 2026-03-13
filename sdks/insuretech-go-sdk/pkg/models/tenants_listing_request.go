package models


// TenantsListingRequest represents a tenants_listing_request
type TenantsListingRequest struct {
	Type string `json:"type"`
	Status string `json:"status,omitempty"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}
