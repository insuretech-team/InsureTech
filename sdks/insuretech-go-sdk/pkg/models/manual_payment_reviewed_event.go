package models

import (
	"time"
)

// ManualPaymentReviewedEvent represents a manual_payment_reviewed_event
type ManualPaymentReviewedEvent struct {
	Approved bool `json:"approved,omitempty"`
	CausationId string `json:"causation_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	RejectionReason string `json:"rejection_reason,omitempty"`
	ReviewNotes string `json:"review_notes,omitempty"`
	ReviewedBy string `json:"reviewed_by,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
}
