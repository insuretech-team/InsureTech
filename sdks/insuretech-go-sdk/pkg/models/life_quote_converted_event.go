package models

import (
	"time"
)

// LifeQuoteConvertedEvent represents a life_quote_converted_event
type LifeQuoteConvertedEvent struct {
	ConvertedAt time.Time `json:"converted_at,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	QuoteNumber string `json:"quote_number,omitempty"`
	TotalPremium string `json:"total_premium,omitempty"`
}
