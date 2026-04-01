package models

import (
	"time"
)

// CommissionPayout represents a commission_payout
type CommissionPayout struct {
	AuditInfo interface{} `json:"audit_info"`
	CommissionCount int `json:"commission_count"`
	Id string `json:"id"`
	PaidAt time.Time `json:"paid_at,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
	PaymentReference string `json:"payment_reference,omitempty"`
	PayoutNumber string `json:"payout_number"`
	PeriodEnd time.Time `json:"period_end"`
	PeriodStart time.Time `json:"period_start"`
	RecipientId string `json:"recipient_id"`
	RecipientType string `json:"recipient_type"`
	Status interface{} `json:"status"`
	TotalAmount *Money `json:"total_amount,omitempty"`
}
