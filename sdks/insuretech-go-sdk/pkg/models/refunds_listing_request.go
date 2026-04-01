package models


// RefundsListingRequest represents a refunds_listing_request
type RefundsListingRequest struct {
	BeneficiaryId string `json:"beneficiary_id"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	Status string `json:"status,omitempty"`
}
