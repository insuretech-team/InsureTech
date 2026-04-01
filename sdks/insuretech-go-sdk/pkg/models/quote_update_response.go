package models


// QuoteUpdateResponse represents a quote_update_response
type QuoteUpdateResponse struct {
	Quote *UnderwritingQuote `json:"quote,omitempty"`
}
