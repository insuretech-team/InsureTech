package models


// InsurersListingRequest represents a insurers_listing_request
type InsurersListingRequest struct {
	Type string `json:"type"`
	Status string `json:"status,omitempty"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}
