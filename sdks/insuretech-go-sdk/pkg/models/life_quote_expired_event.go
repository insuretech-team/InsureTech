package models

import (
	"time"
)

// LifeQuoteExpiredEvent represents a life_quote_expired_event
type LifeQuoteExpiredEvent struct {
	EventId string `json:"event_id,omitempty"`
	ExpiredAt time.Time `json:"expired_at,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	QuoteNumber string `json:"quote_number,omitempty"`
}
