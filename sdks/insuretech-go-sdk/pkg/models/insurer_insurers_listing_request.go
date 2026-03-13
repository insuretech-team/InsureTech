package models


// InsurerInsurersListingRequest represents a insurer_insurers_listing_request
type InsurerInsurersListingRequest struct {
	PageSize int `json:"page_size,omitempty"`
	Type string `json:"type"`
	Status string `json:"status,omitempty"`
	Page int `json:"page,omitempty"`
}
