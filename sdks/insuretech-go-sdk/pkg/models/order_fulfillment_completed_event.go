package models

import (
	"time"
)

// OrderFulfillmentCompletedEvent represents a order_fulfillment_completed_event
type OrderFulfillmentCompletedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	CausationId string `json:"causation_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	InvoiceId string `json:"invoice_id,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
}
