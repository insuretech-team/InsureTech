package models


// UnderwritingQuotesListingRequest represents a underwriting_quotes_listing_request
type UnderwritingQuotesListingRequest struct {
	Status string `json:"status,omitempty"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	BeneficiaryId string `json:"beneficiary_id"`
}
