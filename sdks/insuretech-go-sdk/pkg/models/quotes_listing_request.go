package models


// QuotesListingRequest represents a quotes_listing_request
type QuotesListingRequest struct {
	BeneficiaryId string `json:"beneficiary_id"`
	Status string `json:"status,omitempty"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}
