package models

import (
	"time"
)

// PaymentRefund represents a payment_refund
type PaymentRefund struct {
	ApprovedAt time.Time `json:"approved_at,omitempty"`
	ApprovedBy string `json:"approved_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	PaymentId string `json:"payment_id"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
	Reason string `json:"reason"`
	RefundAmount *Money `json:"refund_amount"`
	RefundId string `json:"refund_id"`
	RefundPaymentId string `json:"refund_payment_id,omitempty"`
	Status interface{} `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}
