package models

import (
	"time"
)

// InvoicePaidEvent represents a invoice_paid_event
type InvoicePaidEvent struct {
	TotalAmount *Money `json:"total_amount,omitempty"`
	CausationId string `json:"causation_id,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	PaidAt time.Time `json:"paid_at,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	EventId string `json:"event_id,omitempty"`
	InvoiceId string `json:"invoice_id,omitempty"`
	InvoiceNumber string `json:"invoice_number,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
}
