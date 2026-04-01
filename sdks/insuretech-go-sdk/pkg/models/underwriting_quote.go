package models

import (
	"time"
)

// UnderwritingQuote represents a underwriting_quote
type UnderwritingQuote struct {
	ApplicantAge int `json:"applicant_age"`
	ApplicantOccupation string `json:"applicant_occupation,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	BasePremium *Money `json:"base_premium,omitempty"`
	BeneficiaryId string `json:"beneficiary_id"`
	ConvertedAt time.Time `json:"converted_at,omitempty"`
	ConvertedPolicyId string `json:"converted_policy_id,omitempty"`
	Id string `json:"id"`
	InsurerProductId string `json:"insurer_product_id"`
	PremiumCalculation string `json:"premium_calculation,omitempty"`
	PremiumPaymentMode string `json:"premium_payment_mode"`
	QuoteNumber string `json:"quote_number"`
	RiderPremium *Money `json:"rider_premium,omitempty"`
	SelectedRiders string `json:"selected_riders,omitempty"`
	Smoker bool `json:"smoker,omitempty"`
	Status interface{} `json:"status"`
	SumAssured *Money `json:"sum_assured,omitempty"`
	TaxAmount *Money `json:"tax_amount,omitempty"`
	TermYears int `json:"term_years"`
	TotalPremium *Money `json:"total_premium,omitempty"`
	ValidUntil time.Time `json:"valid_until"`
}
