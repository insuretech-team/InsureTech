package models

import (
	"time"
)

// ManualPaymentReviewedEvent represents a manual_payment_reviewed_event
type ManualPaymentReviewedEvent struct {
	EventId string `json:"event_id,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	ReviewedBy string `json:"reviewed_by,omitempty"`
	RejectionReason string `json:"rejection_reason,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	Approved bool `json:"approved,omitempty"`
	ReviewNotes string `json:"review_notes,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	CausationId string `json:"causation_id,omitempty"`
}
