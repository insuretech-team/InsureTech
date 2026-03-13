package models


// RefundApprovalRequest represents a refund_approval_request
type RefundApprovalRequest struct {
	RefundId string `json:"refund_id"`
	ApprovedBy string `json:"approved_by,omitempty"`
	Comments string `json:"comments,omitempty"`
}
