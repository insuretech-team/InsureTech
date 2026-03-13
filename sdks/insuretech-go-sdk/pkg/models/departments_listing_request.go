package models


// DepartmentsListingRequest represents a departments_listing_request
type DepartmentsListingRequest struct {
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	BusinessId string `json:"business_id"`
}
