package models


// LifePremiumCalculationRequest represents a life_premium_calculation_request
type LifePremiumCalculationRequest struct {
	AgeAtEntry int `json:"age_at_entry,omitempty"`
	BonusCodes []string `json:"bonus_codes,omitempty"`
	InsuredPerson *InsuredPerson `json:"insured_person,omitempty"`
	PolicyTermYears int `json:"policy_term_years,omitempty"`
	ProductId string `json:"product_id"`
	SumAssured string `json:"sum_assured,omitempty"`
}
