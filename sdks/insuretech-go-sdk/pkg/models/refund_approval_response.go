package models


// RefundApprovalResponse represents a refund_approval_response
type RefundApprovalResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
