package models


// KYCVerificationStartResponse represents a kyc_verification_start_response
type KYCVerificationStartResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
	KycVerificationId string `json:"kyc_verification_id,omitempty"`
}
