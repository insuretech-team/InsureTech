package models

import (
	"time"
)

// InvoicePDFGeneratedEvent represents a invoice_pdfgenerated_event
type InvoicePDFGeneratedEvent struct {
	CausationId string `json:"causation_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	InvoiceId string `json:"invoice_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	InvoiceNumber string `json:"invoice_number,omitempty"`
	InvoicePdfUrl string `json:"invoice_pdf_url,omitempty"`
	DocumentId string `json:"document_id,omitempty"`
}
