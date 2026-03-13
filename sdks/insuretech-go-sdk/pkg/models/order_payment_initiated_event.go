package models

import (
	"time"
)

// OrderPaymentInitiatedEvent represents a order_payment_initiated_event
type OrderPaymentInitiatedEvent struct {
	IdempotencyKey string `json:"idempotency_key,omitempty"`
	PaymentGatewayRef string `json:"payment_gateway_ref,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	Portal string `json:"portal,omitempty"`
	ActorUserId string `json:"actor_user_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	Amount *Money `json:"amount,omitempty"`
	CausationId string `json:"causation_id,omitempty"`
}
