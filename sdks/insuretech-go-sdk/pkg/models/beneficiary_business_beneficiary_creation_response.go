package models


// BeneficiaryBusinessBeneficiaryCreationResponse represents a beneficiary_business_beneficiary_creation_response
type BeneficiaryBusinessBeneficiaryCreationResponse struct {
	BeneficiaryId string `json:"beneficiary_id,omitempty"`
	BeneficiaryCode string `json:"beneficiary_code,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
