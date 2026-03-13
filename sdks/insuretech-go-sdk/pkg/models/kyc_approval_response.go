package models


// KYCApprovalResponse represents a kyc_approval_response
type KYCApprovalResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
