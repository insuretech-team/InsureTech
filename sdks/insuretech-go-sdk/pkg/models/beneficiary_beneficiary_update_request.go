package models


// BeneficiaryBeneficiaryUpdateRequest represents a beneficiary_beneficiary_update_request
type BeneficiaryBeneficiaryUpdateRequest struct {
	Address string `json:"address,omitempty"`
	BeneficiaryId string `json:"beneficiary_id"`
	Email string `json:"email"`
	MobileNumber string `json:"mobile_number,omitempty"`
}
