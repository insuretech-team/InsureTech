package models

import (
	"time"
)

// InvoiceCreatedEvent represents a invoice_created_event
type InvoiceCreatedEvent struct {
	TenantId string `json:"tenant_id,omitempty"`
	DueAt time.Time `json:"due_at,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	CausationId string `json:"causation_id,omitempty"`
	ActorUserId string `json:"actor_user_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	Currency string `json:"currency,omitempty"`
	EventId string `json:"event_id,omitempty"`
	InvoiceId string `json:"invoice_id,omitempty"`
	InvoiceNumber string `json:"invoice_number,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	BusinessId string `json:"business_id,omitempty"`
	Portal string `json:"portal,omitempty"`
	Amount *Money `json:"amount,omitempty"`
	TotalAmount *Money `json:"total_amount,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	PurchaseOrderId string `json:"purchase_order_id,omitempty"`
}
