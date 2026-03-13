package models


// BeneficiaryUpdateRequest represents a beneficiary_update_request
type BeneficiaryUpdateRequest struct {
	BeneficiaryId string `json:"beneficiary_id"`
	MobileNumber string `json:"mobile_number,omitempty"`
	Email string `json:"email"`
	Address string `json:"address,omitempty"`
}
