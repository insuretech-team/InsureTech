package models

import (
	"time"
)

// Policy represents a policy
type Policy struct {
	PaymentGatewayReference string `json:"payment_gateway_reference,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	UnderwritingDecisionId string `json:"underwriting_decision_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	ReceiptNumber string `json:"receipt_number,omitempty"`
	UnderwritingData string `json:"underwriting_data,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Nominees []*Nominee `json:"nominees,omitempty"`
	Riders []*PolicyRider `json:"riders,omitempty"`
	PremiumCurrency string `json:"premium_currency,omitempty"`
	SumInsuredCurrency string `json:"sum_insured_currency,omitempty"`
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate time.Time `json:"end_date,omitempty"`
	VatTax *Money `json:"vat_tax,omitempty"`
	ProviderName string `json:"provider_name,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	PolicyNumber string `json:"policy_number,omitempty"`
	AgentId string `json:"agent_id,omitempty"`
	Status *PolicyStatus `json:"status,omitempty"`
	TenureMonths int `json:"tenure_months,omitempty"`
	PaymentFrequency string `json:"payment_frequency,omitempty"`
	TotalPayable *Money `json:"total_payable,omitempty"`
	OccupationRiskClass string `json:"occupation_risk_class,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	SumInsured *Money `json:"sum_insured,omitempty"`
	IssuedAt time.Time `json:"issued_at,omitempty"`
	PolicyDocumentUrl string `json:"policy_document_url,omitempty"`
	ServiceFee *Money `json:"service_fee,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	PartnerId string `json:"partner_id,omitempty"`
	ProposerDetails *Applicant `json:"proposer_details,omitempty"`
	EnrollmentStartDate time.Time `json:"enrollment_start_date,omitempty"`
	EnrollmentEndDate time.Time `json:"enrollment_end_date,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	PremiumAmount *Money `json:"premium_amount,omitempty"`
	HasExistingPolicies bool `json:"has_existing_policies,omitempty"`
	ClaimsHistorySummary string `json:"claims_history_summary,omitempty"`
}
