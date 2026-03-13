package models


// PremiumCalculationRequest represents a premium_calculation_request
type PremiumCalculationRequest struct {
	ProductId string `json:"product_id"`
	SumInsured *Money `json:"sum_insured,omitempty"`
	TenureMonths int `json:"tenure_months,omitempty"`
	RiderIds []string `json:"rider_ids,omitempty"`
	ApplicantData map[string]interface{} `json:"applicant_data,omitempty"`
}
