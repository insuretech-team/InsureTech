package models


// UnderwritingQuotesListingRequest represents a underwriting_quotes_listing_request
type UnderwritingQuotesListingRequest struct {
	BeneficiaryId string `json:"beneficiary_id"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	Status string `json:"status,omitempty"`
}
