package models


// InsuranceQuotesListingRequest represents a insurance_quotes_listing_request
type InsuranceQuotesListingRequest struct {
	BeneficiaryId string `json:"beneficiary_id"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}
