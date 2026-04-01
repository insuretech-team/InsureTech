package models

import (
	"time"
)

// PaymentRefundProcessedEvent represents a payment_refund_processed_event
type PaymentRefundProcessedEvent struct {
	ActorUserId string `json:"actor_user_id,omitempty"`
	Amount *Money `json:"amount,omitempty"`
	CausationId string `json:"causation_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	InvoiceId string `json:"invoice_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	OriginalPaymentId string `json:"original_payment_id,omitempty"`
	Portal string `json:"portal,omitempty"`
	Reason string `json:"reason,omitempty"`
	RecipientId string `json:"recipient_id,omitempty"`
	RefundId string `json:"refund_id,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
