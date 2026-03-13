package models


// KycKYCRejectionResponse represents a kyc_kyc_rejection_response
type KycKYCRejectionResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
