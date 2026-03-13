package models


// ClaimSubmissionResponse represents a claim_submission_response
type ClaimSubmissionResponse struct {
	ClaimId string `json:"claim_id,omitempty"`
	ClaimNumber string `json:"claim_number,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
