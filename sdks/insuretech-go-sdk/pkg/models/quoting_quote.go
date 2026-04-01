package models

import (
	"time"
)

// QuotingQuote represents a quoting_quote
type QuotingQuote struct {
	AgentId string `json:"agent_id,omitempty"`
	ConvertedAt time.Time `json:"converted_at,omitempty"`
	ConvertedPolicyId string `json:"converted_policy_id,omitempty"`
	CoveragesJson string `json:"coverages_json"`
	CreatedAt time.Time `json:"created_at"`
	CustomerId string `json:"customer_id"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	DiscountsJson string `json:"discounts_json,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	ParametersJson string `json:"parameters_json"`
	ParentQuoteId string `json:"parent_quote_id,omitempty"`
	PremiumCalculationJson string `json:"premium_calculation_json"`
	ProductId string `json:"product_id"`
	QuoteId string `json:"quote_id"`
	QuoteNumber string `json:"quote_number"`
	RevisionNumber int `json:"revision_number"`
	RevisionReason string `json:"revision_reason,omitempty"`
	Status interface{} `json:"status"`
	TotalPremium *Money `json:"total_premium,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
	ValidFrom time.Time `json:"valid_from"`
	ValidUntil time.Time `json:"valid_until"`
}
