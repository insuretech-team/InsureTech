package models

import (
	"time"
)

// PaymentVerifiedEvent represents a payment_verified_event
type PaymentVerifiedEvent struct {
	PaymentId string `json:"payment_id,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	VerifiedBy string `json:"verified_by,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	ValId string `json:"val_id,omitempty"`
	Amount *Money `json:"amount,omitempty"`
	VerifiedAt time.Time `json:"verified_at,omitempty"`
	EventId string `json:"event_id,omitempty"`
	InvoiceId string `json:"invoice_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	CausationId string `json:"causation_id,omitempty"`
	Provider string `json:"provider,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
}
