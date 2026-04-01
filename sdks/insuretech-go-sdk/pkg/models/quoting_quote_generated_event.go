package models

import (
	"time"
)

// QuotingQuoteGeneratedEvent represents a quoting_quote_generated_event
type QuotingQuoteGeneratedEvent struct {
	AgentId string `json:"agent_id,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	GeneratedAt time.Time `json:"generated_at,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	QuoteNumber string `json:"quote_number,omitempty"`
	TotalPremium *Money `json:"total_premium,omitempty"`
	ValidUntil time.Time `json:"valid_until,omitempty"`
}
