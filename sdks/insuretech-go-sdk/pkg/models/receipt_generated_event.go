package models

import (
	"time"
)

// ReceiptGeneratedEvent represents a receipt_generated_event
type ReceiptGeneratedEvent struct {
	CausationId string `json:"causation_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	ReceiptDocumentId string `json:"receipt_document_id,omitempty"`
	ReceiptFileId string `json:"receipt_file_id,omitempty"`
	ReceiptNumber string `json:"receipt_number,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
}
