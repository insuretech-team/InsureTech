package models

import (
	"time"
)

// PaymentFailedEvent represents a payment_failed_event
type PaymentFailedEvent struct {
	PayerId string `json:"payer_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	ActorUserId string `json:"actor_user_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	Amount *Money `json:"amount,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	Portal string `json:"portal,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	InvoiceId string `json:"invoice_id,omitempty"`
	CausationId string `json:"causation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	Provider string `json:"provider,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
}
