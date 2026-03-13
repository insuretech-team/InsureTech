package models


// InitiateRefundRequest represents a initiate_refund_request
type InitiateRefundRequest struct {
	Reason string `json:"reason,omitempty"`
	InitiatedBy string `json:"initiated_by,omitempty"`
	PaymentId string `json:"payment_id"`
	RefundAmount *Money `json:"refund_amount,omitempty"`
}
