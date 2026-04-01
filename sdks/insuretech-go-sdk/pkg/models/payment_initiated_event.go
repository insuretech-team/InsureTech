package models

import (
	"time"
)

// PaymentInitiatedEvent represents a payment_initiated_event
type PaymentInitiatedEvent struct {
	ActorUserId string `json:"actor_user_id,omitempty"`
	Amount *Money `json:"amount,omitempty"`
	CausationId string `json:"causation_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	InvoiceId string `json:"invoice_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	PayerId string `json:"payer_id,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
	Portal string `json:"portal,omitempty"`
	Provider string `json:"provider,omitempty"`
	ReferenceId string `json:"reference_id,omitempty"`
	ReferenceType string `json:"reference_type,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	TranId string `json:"tran_id,omitempty"`
}
