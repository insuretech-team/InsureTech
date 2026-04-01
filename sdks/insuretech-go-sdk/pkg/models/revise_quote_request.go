package models


// ReviseQuoteRequest represents a revise_quote_request
type ReviseQuoteRequest struct {
	NewParameters *QuoteParameters `json:"new_parameters,omitempty"`
	QuoteId string `json:"quote_id"`
	RevisionReason string `json:"revision_reason,omitempty"`
	ValidityDays int `json:"validity_days,omitempty"`
}
