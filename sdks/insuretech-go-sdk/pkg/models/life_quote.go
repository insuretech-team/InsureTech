package models

import (
	"time"
)

// LifeQuote represents a life_quote
type LifeQuote struct {
	AgeAddition string `json:"age_addition"`
	AgeAtEntry int `json:"age_at_entry"`
	AgentId string `json:"agent_id,omitempty"`
	BasePremium string `json:"base_premium"`
	BonusDiscount string `json:"bonus_discount"`
	BonusesAppliedJson string `json:"bonuses_applied_json,omitempty"`
	ConditionAddition string `json:"condition_addition"`
	ConditionMultiplier float64 `json:"condition_multiplier"`
	ConvertedAt time.Time `json:"converted_at,omitempty"`
	ConvertedPolicyId string `json:"converted_policy_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	CustomerId string `json:"customer_id"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	HealthConditionsJson string `json:"health_conditions_json,omitempty"`
	InsuredPersonJson string `json:"insured_person_json"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	PolicyTermYears int `json:"policy_term_years"`
	ProductId string `json:"product_id"`
	QuoteId string `json:"quote_id"`
	QuoteNumber string `json:"quote_number"`
	Status interface{} `json:"status"`
	SumAssured string `json:"sum_assured"`
	TotalPremium string `json:"total_premium"`
	UpdatedAt time.Time `json:"updated_at"`
	ValidUntil time.Time `json:"valid_until"`
}
