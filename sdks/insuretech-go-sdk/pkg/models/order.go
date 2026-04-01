package models

import (
	"time"
)

// Order represents a order
type Order struct {
	ActorUserId string `json:"actor_user_id,omitempty"`
	BillingStatus interface{} `json:"billing_status"`
	CancellationReason string `json:"cancellation_reason,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	CoverageEndAt time.Time `json:"coverage_end_at,omitempty"`
	CoverageStartAt time.Time `json:"coverage_start_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	Currency string `json:"currency"`
	CustomerId string `json:"customer_id"`
	FailureReason string `json:"failure_reason,omitempty"`
	FulfillmentStatus interface{} `json:"fulfillment_status"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
	InsurerId string `json:"insurer_id,omitempty"`
	InvoiceId string `json:"invoice_id,omitempty"`
	ManualReviewRequired bool `json:"manual_review_required"`
	OrderId string `json:"order_id"`
	OrderNumber string `json:"order_number"`
	OrganisationId string `json:"organisation_id,omitempty"`
	PaidAt time.Time `json:"paid_at,omitempty"`
	PaymentDueAt time.Time `json:"payment_due_at,omitempty"`
	PaymentGatewayRef string `json:"payment_gateway_ref,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	PaymentStatus interface{} `json:"payment_status"`
	PlanId string `json:"plan_id"`
	PolicyId string `json:"policy_id,omitempty"`
	Portal string `json:"portal,omitempty"`
	ProductId string `json:"product_id"`
	ProposalDecidedAt time.Time `json:"proposal_decided_at,omitempty"`
	ProposalDecisionReason string `json:"proposal_decision_reason,omitempty"`
	ProposalId string `json:"proposal_id,omitempty"`
	ProposalStatus interface{} `json:"proposal_status"`
	ProposalSubmittedAt time.Time `json:"proposal_submitted_at,omitempty"`
	PurchaseOrderId string `json:"purchase_order_id,omitempty"`
	QuotationId string `json:"quotation_id"`
	RefundId string `json:"refund_id,omitempty"`
	Status interface{} `json:"status"`
	TenantId string `json:"tenant_id"`
	TotalPayable *Money `json:"total_payable"`
	UpdatedAt time.Time `json:"updated_at"`
}
