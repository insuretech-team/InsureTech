package models

import (
	"time"
)

// QuotingQuoteConvertedEvent represents a quoting_quote_converted_event
type QuotingQuoteConvertedEvent struct {
	ConvertedAt time.Time `json:"converted_at,omitempty"`
	ConvertedPremium *Money `json:"converted_premium,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	QuoteNumber string `json:"quote_number,omitempty"`
}
