package models

import (
	"time"
)

// OrderCreationRequest represents a order_creation_request
type OrderCreationRequest struct {
	CoverageEndAt time.Time `json:"coverage_end_at,omitempty"`
	CoverageStartAt time.Time `json:"coverage_start_at,omitempty"`
	CustomerId string `json:"customer_id"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
	OrganisationId string `json:"organisation_id"`
	PaymentDueAt time.Time `json:"payment_due_at,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
	PlanId string `json:"plan_id"`
	ProductId string `json:"product_id"`
	PurchaseOrderId string `json:"purchase_order_id"`
	QuotationId string `json:"quotation_id"`
	TenantId string `json:"tenant_id"`
	TotalPayable *Money `json:"total_payable,omitempty"`
}
