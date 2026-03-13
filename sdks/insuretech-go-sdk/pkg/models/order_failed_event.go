package models

import (
	"time"
)

// OrderFailedEvent represents a order_failed_event
type OrderFailedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	Portal string `json:"portal,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	CausationId string `json:"causation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	ActorUserId string `json:"actor_user_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}
