package models

import (
	"time"
)

// OrderCreatedEvent represents a order_created_event
type OrderCreatedEvent struct {
	ActorUserId string `json:"actor_user_id,omitempty"`
	CausationId string `json:"causation_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	OrderNumber string `json:"order_number,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	PlanId string `json:"plan_id,omitempty"`
	Portal string `json:"portal,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	PurchaseOrderId string `json:"purchase_order_id,omitempty"`
	QuotationId string `json:"quotation_id,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	SessionType string `json:"session_type,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	TotalPayable *Money `json:"total_payable,omitempty"`
}
