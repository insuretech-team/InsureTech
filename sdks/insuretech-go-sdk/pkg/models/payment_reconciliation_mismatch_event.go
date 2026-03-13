package models

import (
	"time"
)

// PaymentReconciliationMismatchEvent represents a payment_reconciliation_mismatch_event
type PaymentReconciliationMismatchEvent struct {
	EventId string `json:"event_id,omitempty"`
	ExpectedAmount string `json:"expected_amount,omitempty"`
	ActualAmount string `json:"actual_amount,omitempty"`
	GatewayRef string `json:"gateway_ref,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	ReconciliationId string `json:"reconciliation_id,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	Provider string `json:"provider,omitempty"`
	MismatchReason string `json:"mismatch_reason,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
}
