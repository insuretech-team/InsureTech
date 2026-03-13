package models


// BeneficiaryBeneficiaryUpdateRequest represents a beneficiary_beneficiary_update_request
type BeneficiaryBeneficiaryUpdateRequest struct {
	MobileNumber string `json:"mobile_number,omitempty"`
	Email string `json:"email"`
	Address string `json:"address,omitempty"`
	BeneficiaryId string `json:"beneficiary_id"`
}
