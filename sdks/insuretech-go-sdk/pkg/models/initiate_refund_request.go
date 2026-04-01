package models


// InitiateRefundRequest represents a initiate_refund_request
type InitiateRefundRequest struct {
	InitiatedBy string `json:"initiated_by,omitempty"`
	PaymentId string `json:"payment_id"`
	Reason string `json:"reason,omitempty"`
	RefundAmount *Money `json:"refund_amount,omitempty"`
}
