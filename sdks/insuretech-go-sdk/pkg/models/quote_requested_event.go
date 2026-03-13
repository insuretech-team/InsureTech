package models

import (
	"time"
)

// QuoteRequestedEvent represents a quote_requested_event
type QuoteRequestedEvent struct {
	SumAssured *Money `json:"sum_assured,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	QuoteNumber string `json:"quote_number,omitempty"`
	BeneficiaryId string `json:"beneficiary_id,omitempty"`
	InsurerProductId string `json:"insurer_product_id,omitempty"`
}
