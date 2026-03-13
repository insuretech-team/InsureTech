package models

import (
	"time"
)

// InvoiceCancelledEvent represents a invoice_cancelled_event
type InvoiceCancelledEvent struct {
	EventId string `json:"event_id,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	CancelledBy string `json:"cancelled_by,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	CreditNoteId string `json:"credit_note_id,omitempty"`
	InvoiceId string `json:"invoice_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	CreditNoteIssued bool `json:"credit_note_issued,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	InvoiceNumber string `json:"invoice_number,omitempty"`
	Reason string `json:"reason,omitempty"`
	CausationId string `json:"causation_id,omitempty"`
	CancelledAt time.Time `json:"cancelled_at,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
}
