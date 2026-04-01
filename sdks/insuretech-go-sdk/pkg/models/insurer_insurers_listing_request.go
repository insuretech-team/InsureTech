package models


// InsurerInsurersListingRequest represents a insurer_insurers_listing_request
type InsurerInsurersListingRequest struct {
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	Status string `json:"status,omitempty"`
	Type string `json:"type"`
}
