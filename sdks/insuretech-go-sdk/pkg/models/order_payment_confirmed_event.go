package models

import (
	"time"
)

// OrderPaymentConfirmedEvent represents a order_payment_confirmed_event
type OrderPaymentConfirmedEvent struct {
	OrderId string `json:"order_id,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	PlanId string `json:"plan_id,omitempty"`
	Portal string `json:"portal,omitempty"`
	TotalPayable *Money `json:"total_payable,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	ActorUserId string `json:"actor_user_id,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	CausationId string `json:"causation_id,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	QuotationId string `json:"quotation_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	TenantId string `json:"tenant_id,omitempty"`
	InvoiceId string `json:"invoice_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	EventId string `json:"event_id,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
