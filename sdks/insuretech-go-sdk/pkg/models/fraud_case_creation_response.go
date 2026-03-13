package models


// FraudCaseCreationResponse represents a fraud_case_creation_response
type FraudCaseCreationResponse struct {
	CaseNumber string `json:"case_number,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
	FraudCaseId string `json:"fraud_case_id,omitempty"`
}
