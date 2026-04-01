package models


// LifeQuoteGenerationRequest represents a life_quote_generation_request
type LifeQuoteGenerationRequest struct {
	AgeAtEntry int `json:"age_at_entry,omitempty"`
	AgentId string `json:"agent_id"`
	BonusCodes []string `json:"bonus_codes,omitempty"`
	CustomerId string `json:"customer_id"`
	InsuredPerson *InsuredPerson `json:"insured_person,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	PolicyTermYears int `json:"policy_term_years,omitempty"`
	ProductId string `json:"product_id"`
	SumAssured string `json:"sum_assured,omitempty"`
	ValidityDays int `json:"validity_days,omitempty"`
}
