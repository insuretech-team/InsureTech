package models


// BeneficiaryBeneficiariesListingRequest represents a beneficiary_beneficiaries_listing_request
type BeneficiaryBeneficiariesListingRequest struct {
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	Status string `json:"status,omitempty"`
	Type string `json:"type"`
}
