package models

import (
	"time"
)

// Order represents a order
type Order struct {
	Portal string `json:"portal,omitempty"`
	TenantId string `json:"tenant_id"`
	CancellationReason string `json:"cancellation_reason,omitempty"`
	FailureReason string `json:"failure_reason,omitempty"`
	BillingStatus interface{} `json:"billing_status"`
	PlanId string `json:"plan_id"`
	Status interface{} `json:"status"`
	TotalPayable *Money `json:"total_payable"`
	OrganisationId string `json:"organisation_id,omitempty"`
	CoverageStartAt time.Time `json:"coverage_start_at,omitempty"`
	ActorUserId string `json:"actor_user_id,omitempty"`
	PaymentGatewayRef string `json:"payment_gateway_ref,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	PaidAt time.Time `json:"paid_at,omitempty"`
	CustomerId string `json:"customer_id"`
	Currency string `json:"currency"`
	UpdatedAt time.Time `json:"updated_at"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
	ManualReviewRequired bool `json:"manual_review_required"`
	PaymentId string `json:"payment_id,omitempty"`
	InvoiceId string `json:"invoice_id,omitempty"`
	CoverageEndAt time.Time `json:"coverage_end_at,omitempty"`
	PurchaseOrderId string `json:"purchase_order_id,omitempty"`
	OrderNumber string `json:"order_number"`
	ProductId string `json:"product_id"`
	PolicyId string `json:"policy_id,omitempty"`
	PaymentStatus interface{} `json:"payment_status"`
	PaymentDueAt time.Time `json:"payment_due_at,omitempty"`
	OrderId string `json:"order_id"`
	CorrelationId string `json:"correlation_id,omitempty"`
	FulfillmentStatus interface{} `json:"fulfillment_status"`
	QuotationId string `json:"quotation_id"`
}
