package models


// IndividualBeneficiaryCreationResponse represents a individual_beneficiary_creation_response
type IndividualBeneficiaryCreationResponse struct {
	BeneficiaryId string `json:"beneficiary_id,omitempty"`
	BeneficiaryCode string `json:"beneficiary_code,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
