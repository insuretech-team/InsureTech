package models

import (
	"time"
)

// InsuranceProposal represents a insurance_proposal
type InsuranceProposal struct {
	ApprovedPolicyId string `json:"approved_policy_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	CustomerId string `json:"customer_id"`
	DecisionReason string `json:"decision_reason,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	InsurerId string `json:"insurer_id"`
	InsurerResponsePayload string `json:"insurer_response_payload,omitempty"`
	OrderId string `json:"order_id"`
	PlanId string `json:"plan_id"`
	ProductId string `json:"product_id"`
	ProposalId string `json:"proposal_id"`
	ProposalNumber string `json:"proposal_number"`
	ProposedPremium *Money `json:"proposed_premium"`
	ProposedPremiumCurrency string `json:"proposed_premium_currency"`
	ProposedSumInsured *Money `json:"proposed_sum_insured"`
	ProposedSumInsuredCurrency string `json:"proposed_sum_insured_currency"`
	QuotationId string `json:"quotation_id"`
	RefundId string `json:"refund_id,omitempty"`
	ReviewedAt time.Time `json:"reviewed_at,omitempty"`
	ReviewedByUserId string `json:"reviewed_by_user_id,omitempty"`
	Status interface{} `json:"status"`
	SubmissionPayload string `json:"submission_payload,omitempty"`
	SubmittedAt time.Time `json:"submitted_at"`
	TenantId string `json:"tenant_id"`
	UpdatedAt time.Time `json:"updated_at"`
}
