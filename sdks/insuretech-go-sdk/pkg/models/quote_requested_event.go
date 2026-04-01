package models

import (
	"time"
)

// QuoteRequestedEvent represents a quote_requested_event
type QuoteRequestedEvent struct {
	BeneficiaryId string `json:"beneficiary_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	InsurerProductId string `json:"insurer_product_id,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	QuoteNumber string `json:"quote_number,omitempty"`
	SumAssured *Money `json:"sum_assured,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
