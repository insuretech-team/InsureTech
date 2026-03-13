package models


// ClaimApprovalResponse represents a claim_approval_response
type ClaimApprovalResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
