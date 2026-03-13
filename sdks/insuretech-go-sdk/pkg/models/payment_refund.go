package models

import (
	"time"
)

// PaymentRefund represents a payment_refund
type PaymentRefund struct {
	RefundId string `json:"refund_id"`
	RefundPaymentId string `json:"refund_payment_id,omitempty"`
	RefundAmount *Money `json:"refund_amount"`
	Reason string `json:"reason"`
	Status interface{} `json:"status"`
	ApprovedAt time.Time `json:"approved_at,omitempty"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	PaymentId string `json:"payment_id"`
	ApprovedBy string `json:"approved_by,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}
