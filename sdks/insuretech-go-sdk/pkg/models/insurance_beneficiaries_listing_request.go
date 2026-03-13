package models


// InsuranceBeneficiariesListingRequest represents a insurance_beneficiaries_listing_request
type InsuranceBeneficiariesListingRequest struct {
	Page int `json:"page"`
	PageSize int `json:"page_size,omitempty"`
}
