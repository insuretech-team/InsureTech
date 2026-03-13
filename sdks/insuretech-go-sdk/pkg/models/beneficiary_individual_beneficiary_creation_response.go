package models


// BeneficiaryIndividualBeneficiaryCreationResponse represents a beneficiary_individual_beneficiary_creation_response
type BeneficiaryIndividualBeneficiaryCreationResponse struct {
	BeneficiaryId string `json:"beneficiary_id,omitempty"`
	BeneficiaryCode string `json:"beneficiary_code,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
