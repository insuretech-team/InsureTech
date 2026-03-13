package models


// RequestQuoteRequest represents a request_quote_request
type RequestQuoteRequest struct {
	PremiumPaymentMode string `json:"premium_payment_mode,omitempty"`
	RiderCodes []string `json:"rider_codes,omitempty"`
	ApplicantAge int `json:"applicant_age,omitempty"`
	Smoker bool `json:"smoker,omitempty"`
	BeneficiaryId string `json:"beneficiary_id"`
	InsurerProductId string `json:"insurer_product_id"`
	SumAssured *Money `json:"sum_assured,omitempty"`
	TermYears int `json:"term_years,omitempty"`
}
