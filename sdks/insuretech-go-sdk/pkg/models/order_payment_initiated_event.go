package models

import (
	"time"
)

// OrderPaymentInitiatedEvent represents a order_payment_initiated_event
type OrderPaymentInitiatedEvent struct {
	ActorUserId string `json:"actor_user_id,omitempty"`
	Amount *Money `json:"amount,omitempty"`
	CausationId string `json:"causation_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	PaymentGatewayRef string `json:"payment_gateway_ref,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
	Portal string `json:"portal,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
