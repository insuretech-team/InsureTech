package models

import (
	"time"
)

// OrderCancelledEvent represents a order_cancelled_event
type OrderCancelledEvent struct {
	ActorUserId string `json:"actor_user_id,omitempty"`
	CausationId string `json:"causation_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	Portal string `json:"portal,omitempty"`
	Reason string `json:"reason,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
