package models


// RequestQuoteRequest represents a request_quote_request
type RequestQuoteRequest struct {
	ApplicantAge int `json:"applicant_age,omitempty"`
	BeneficiaryId string `json:"beneficiary_id"`
	InsurerProductId string `json:"insurer_product_id"`
	PremiumPaymentMode string `json:"premium_payment_mode,omitempty"`
	RiderCodes []string `json:"rider_codes,omitempty"`
	Smoker bool `json:"smoker,omitempty"`
	SumAssured *Money `json:"sum_assured,omitempty"`
	TermYears int `json:"term_years,omitempty"`
}
