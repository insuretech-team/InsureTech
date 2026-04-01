package models

import (
	"time"
)

// InvoiceCreatedEvent represents a invoice_created_event
type InvoiceCreatedEvent struct {
	ActorUserId string `json:"actor_user_id,omitempty"`
	Amount *Money `json:"amount,omitempty"`
	BusinessId string `json:"business_id,omitempty"`
	CausationId string `json:"causation_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Currency string `json:"currency,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	DueAt time.Time `json:"due_at,omitempty"`
	EventId string `json:"event_id,omitempty"`
	InvoiceId string `json:"invoice_id,omitempty"`
	InvoiceNumber string `json:"invoice_number,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	Portal string `json:"portal,omitempty"`
	PurchaseOrderId string `json:"purchase_order_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	TotalAmount *Money `json:"total_amount,omitempty"`
}
