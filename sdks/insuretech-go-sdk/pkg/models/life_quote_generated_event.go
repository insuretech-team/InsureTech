package models

import (
	"time"
)

// LifeQuoteGeneratedEvent represents a life_quote_generated_event
type LifeQuoteGeneratedEvent struct {
	AgeAtEntry int `json:"age_at_entry,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	GeneratedAt time.Time `json:"generated_at,omitempty"`
	InsuredPersonName string `json:"insured_person_name,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	QuoteNumber string `json:"quote_number,omitempty"`
	SumAssured string `json:"sum_assured,omitempty"`
	TotalPremium string `json:"total_premium,omitempty"`
	ValidUntil time.Time `json:"valid_until,omitempty"`
}
