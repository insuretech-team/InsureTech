package models

import (
	"time"
)

// QuoteRevisedEvent represents a quote_revised_event
type QuoteRevisedEvent struct {
	EventId string `json:"event_id,omitempty"`
	NewPremium *Money `json:"new_premium,omitempty"`
	ParentQuoteId string `json:"parent_quote_id,omitempty"`
	PreviousPremium *Money `json:"previous_premium,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	QuoteNumber string `json:"quote_number,omitempty"`
	RevisedAt time.Time `json:"revised_at,omitempty"`
	RevisionNumber int `json:"revision_number,omitempty"`
	RevisionReason string `json:"revision_reason,omitempty"`
}
