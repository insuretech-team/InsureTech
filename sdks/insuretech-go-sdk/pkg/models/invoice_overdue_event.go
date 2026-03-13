package models

import (
	"time"
)

// InvoiceOverdueEvent represents a invoice_overdue_event
type InvoiceOverdueEvent struct {
	InvoiceNumber string `json:"invoice_number,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	DueAt time.Time `json:"due_at,omitempty"`
	TotalAmount *Money `json:"total_amount,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	EventId string `json:"event_id,omitempty"`
	InvoiceId string `json:"invoice_id,omitempty"`
}
