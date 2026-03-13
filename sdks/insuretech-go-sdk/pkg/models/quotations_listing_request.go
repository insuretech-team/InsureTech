package models


// QuotationsListingRequest represents a quotations_listing_request
type QuotationsListingRequest struct {
	BusinessId string `json:"business_id"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}
