package models


// FraudCaseUpdateResponse represents a fraud_case_update_response
type FraudCaseUpdateResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
