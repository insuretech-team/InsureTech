package models

import (
	"time"
)

// PaymentCompletedEvent represents a payment_completed_event
type PaymentCompletedEvent struct {
	OrderId string `json:"order_id,omitempty"`
	Provider string `json:"provider,omitempty"`
	ValId string `json:"val_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	ActorUserId string `json:"actor_user_id,omitempty"`
	CausationId string `json:"causation_id,omitempty"`
	ReceiptNumber string `json:"receipt_number,omitempty"`
	TransactionId string `json:"transaction_id,omitempty"`
	Amount *Money `json:"amount,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
	InvoiceId string `json:"invoice_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	Portal string `json:"portal,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
