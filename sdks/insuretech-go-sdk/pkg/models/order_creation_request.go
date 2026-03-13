package models

import (
	"time"
)

// OrderCreationRequest represents a order_creation_request
type OrderCreationRequest struct {
	CustomerId string `json:"customer_id"`
	PaymentMethod string `json:"payment_method,omitempty"`
	TenantId string `json:"tenant_id"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
	TotalPayable *Money `json:"total_payable,omitempty"`
	PurchaseOrderId string `json:"purchase_order_id"`
	QuotationId string `json:"quotation_id"`
	OrganisationId string `json:"organisation_id"`
	CoverageStartAt time.Time `json:"coverage_start_at,omitempty"`
	CoverageEndAt time.Time `json:"coverage_end_at,omitempty"`
	PaymentDueAt time.Time `json:"payment_due_at,omitempty"`
	ProductId string `json:"product_id"`
	PlanId string `json:"plan_id"`
}
