package models

import (
	"time"
)

// Refund represents a refund
type Refund struct {
	ApprovedBy string `json:"approved_by,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	CalculationDetails string `json:"calculation_details,omitempty"`
	CancellationCharge *Money `json:"cancellation_charge,omitempty"`
	Id string `json:"id"`
	OrderId string `json:"order_id,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
	PaymentReference string `json:"payment_reference,omitempty"`
	PaymentRefundId string `json:"payment_refund_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	PremiumUsed *Money `json:"premium_used,omitempty"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
	ProposalId string `json:"proposal_id,omitempty"`
	Reason *RefundReason `json:"reason"`
	ReasonDetails string `json:"reason_details,omitempty"`
	RefundNumber string `json:"refund_number"`
	RefundableAmount *Money `json:"refundable_amount,omitempty"`
	RequestedBy string `json:"requested_by"`
	Status interface{} `json:"status"`
	TotalPremiumPaid *Money `json:"total_premium_paid,omitempty"`
}
