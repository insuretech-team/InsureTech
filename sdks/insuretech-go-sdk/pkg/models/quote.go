package models

import (
	"time"
)

// Quote represents a quote
type Quote struct {
	SumAssured *Money `json:"sum_assured,omitempty"`
	TermYears int `json:"term_years"`
	PremiumPaymentMode string `json:"premium_payment_mode"`
	BasePremium *Money `json:"base_premium,omitempty"`
	TaxAmount *Money `json:"tax_amount,omitempty"`
	ConvertedPolicyId string `json:"converted_policy_id,omitempty"`
	BeneficiaryId string `json:"beneficiary_id"`
	RiderPremium *Money `json:"rider_premium,omitempty"`
	TotalPremium *Money `json:"total_premium,omitempty"`
	PremiumCalculation string `json:"premium_calculation,omitempty"`
	ApplicantAge int `json:"applicant_age"`
	Smoker bool `json:"smoker,omitempty"`
	ValidUntil time.Time `json:"valid_until"`
	AuditInfo interface{} `json:"audit_info"`
	Status interface{} `json:"status"`
	ConvertedAt time.Time `json:"converted_at,omitempty"`
	Id string `json:"id"`
	QuoteNumber string `json:"quote_number"`
	InsurerProductId string `json:"insurer_product_id"`
	SelectedRiders string `json:"selected_riders,omitempty"`
	ApplicantOccupation string `json:"applicant_occupation,omitempty"`
}
