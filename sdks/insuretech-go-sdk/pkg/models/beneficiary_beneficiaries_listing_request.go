package models


// BeneficiaryBeneficiariesListingRequest represents a beneficiary_beneficiaries_listing_request
type BeneficiaryBeneficiariesListingRequest struct {
	Type string `json:"type"`
	Status string `json:"status,omitempty"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}
