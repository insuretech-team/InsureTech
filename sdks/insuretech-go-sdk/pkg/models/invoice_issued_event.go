package models

import (
	"time"
)

// InvoiceIssuedEvent represents a invoice_issued_event
type InvoiceIssuedEvent struct {
	CausationId string `json:"causation_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	InvoiceId string `json:"invoice_id,omitempty"`
	InvoiceNumber string `json:"invoice_number,omitempty"`
	IssuedAt time.Time `json:"issued_at,omitempty"`
	IssuedBy string `json:"issued_by,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	TotalAmount *Money `json:"total_amount,omitempty"`
}
