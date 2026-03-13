package models


// BeneficiaryUpdateResponse represents a beneficiary_update_response
type BeneficiaryUpdateResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
