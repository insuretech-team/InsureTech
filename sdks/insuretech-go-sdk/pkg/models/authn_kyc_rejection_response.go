package models


// AuthnKYCRejectionResponse represents a authn_kyc_rejection_response
type AuthnKYCRejectionResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
