package models


// KYCRejectionResponse represents a kyc_rejection_response
type KYCRejectionResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
