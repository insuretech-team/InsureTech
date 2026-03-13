package models

import (
	"time"
)

// Refund represents a refund
type Refund struct {
	PremiumUsed *Money `json:"premium_used,omitempty"`
	Status interface{} `json:"status"`
	ApprovedBy string `json:"approved_by,omitempty"`
	RefundNumber string `json:"refund_number"`
	PolicyId string `json:"policy_id"`
	PaymentMethod string `json:"payment_method,omitempty"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	Id string `json:"id"`
	Reason *RefundReason `json:"reason"`
	CalculationDetails string `json:"calculation_details,omitempty"`
	RequestedBy string `json:"requested_by"`
	PaymentReference string `json:"payment_reference,omitempty"`
	PaymentRefundId string `json:"payment_refund_id,omitempty"`
	ReasonDetails string `json:"reason_details,omitempty"`
	TotalPremiumPaid *Money `json:"total_premium_paid,omitempty"`
	CancellationCharge *Money `json:"cancellation_charge,omitempty"`
	RefundableAmount *Money `json:"refundable_amount,omitempty"`
}
