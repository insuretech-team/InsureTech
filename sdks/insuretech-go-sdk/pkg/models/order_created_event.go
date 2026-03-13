package models

import (
	"time"
)

// OrderCreatedEvent represents a order_created_event
type OrderCreatedEvent struct {
	QuotationId string `json:"quotation_id,omitempty"`
	CausationId string `json:"causation_id,omitempty"`
	OrderId string `json:"order_id,omitempty"`
	OrderNumber string `json:"order_number,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	PlanId string `json:"plan_id,omitempty"`
	TotalPayable *Money `json:"total_payable,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	Portal string `json:"portal,omitempty"`
	SessionType string `json:"session_type,omitempty"`
	PurchaseOrderId string `json:"purchase_order_id,omitempty"`
	ActorUserId string `json:"actor_user_id,omitempty"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
}
