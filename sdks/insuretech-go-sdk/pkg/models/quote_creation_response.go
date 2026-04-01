package models


// QuoteCreationResponse represents a quote_creation_response
type QuoteCreationResponse struct {
	Quote *UnderwritingQuote `json:"quote,omitempty"`
}
