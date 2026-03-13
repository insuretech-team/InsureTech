package models

import (
	"time"
)

// ReceiptGeneratedEvent represents a receipt_generated_event
type ReceiptGeneratedEvent struct {
	EventId string `json:"event_id,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	ReceiptNumber string `json:"receipt_number,omitempty"`
	ReceiptDocumentId string `json:"receipt_document_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	CausationId string `json:"causation_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	ReceiptFileId string `json:"receipt_file_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
}
