package models


// KYCCompletionResponse represents a kyc_completion_response
type KYCCompletionResponse struct {
	KycStatus string `json:"kyc_status,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
