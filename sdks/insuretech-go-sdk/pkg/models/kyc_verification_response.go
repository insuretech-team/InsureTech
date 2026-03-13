package models


// KYCVerificationResponse represents a kyc_verification_response
type KYCVerificationResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
