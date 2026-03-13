package models


// ClaimRejectionResponse represents a claim_rejection_response
type ClaimRejectionResponse struct {
	Error *Error `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}
