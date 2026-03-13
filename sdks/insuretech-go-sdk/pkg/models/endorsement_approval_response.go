package models


// EndorsementApprovalResponse represents a endorsement_approval_response
type EndorsementApprovalResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
