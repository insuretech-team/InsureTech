package models

import (
	"time"
)

// Policy represents a policy
type Policy struct {
	AgentId string `json:"agent_id,omitempty"`
	ClaimsHistorySummary string `json:"claims_history_summary,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	EndDate time.Time `json:"end_date,omitempty"`
	EnrollmentEndDate time.Time `json:"enrollment_end_date,omitempty"`
	EnrollmentStartDate time.Time `json:"enrollment_start_date,omitempty"`
	HasExistingPolicies bool `json:"has_existing_policies,omitempty"`
	IssuedAt time.Time `json:"issued_at,omitempty"`
	Nominees []*Nominee `json:"nominees,omitempty"`
	OccupationRiskClass string `json:"occupation_risk_class,omitempty"`
	PartnerId string `json:"partner_id,omitempty"`
	PaymentFrequency string `json:"payment_frequency,omitempty"`
	PaymentGatewayReference string `json:"payment_gateway_reference,omitempty"`
	PolicyDocumentUrl string `json:"policy_document_url,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	PolicyNumber string `json:"policy_number,omitempty"`
	PremiumAmount *Money `json:"premium_amount,omitempty"`
	PremiumCurrency string `json:"premium_currency,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	ProposerDetails *Applicant `json:"proposer_details,omitempty"`
	ProviderName string `json:"provider_name,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	ReceiptNumber string `json:"receipt_number,omitempty"`
	Riders []*PolicyRider `json:"riders,omitempty"`
	ServiceFee *Money `json:"service_fee,omitempty"`
	StartDate time.Time `json:"start_date,omitempty"`
	Status *PolicyStatus `json:"status,omitempty"`
	SumInsured *Money `json:"sum_insured,omitempty"`
	SumInsuredCurrency string `json:"sum_insured_currency,omitempty"`
	TenureMonths int `json:"tenure_months,omitempty"`
	TotalPayable *Money `json:"total_payable,omitempty"`
	UnderwritingData string `json:"underwriting_data,omitempty"`
	UnderwritingDecisionId string `json:"underwriting_decision_id,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	VatTax *Money `json:"vat_tax,omitempty"`
}
