package models


// ProductsPremiumCalculationRequest represents a products_premium_calculation_request
type ProductsPremiumCalculationRequest struct {
	ApplicantData map[string]interface{} `json:"applicant_data,omitempty"`
	ProductId string `json:"product_id"`
	RiderIds []string `json:"rider_ids,omitempty"`
	SumInsured *Money `json:"sum_insured,omitempty"`
	TenureMonths int `json:"tenure_months,omitempty"`
}
