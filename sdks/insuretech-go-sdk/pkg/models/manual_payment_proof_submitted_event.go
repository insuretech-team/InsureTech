package models

import (
	"time"
)

// ManualPaymentProofSubmittedEvent represents a manual_payment_proof_submitted_event
type ManualPaymentProofSubmittedEvent struct {
	SubmittedBy string `json:"submitted_by,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	CausationId string `json:"causation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	ManualProofFileId string `json:"manual_proof_file_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
}
