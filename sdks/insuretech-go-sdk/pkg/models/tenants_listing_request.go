package models


// TenantsListingRequest represents a tenants_listing_request
type TenantsListingRequest struct {
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	Status string `json:"status,omitempty"`
	Type string `json:"type"`
}
