package models


// ReviseQuoteResponse represents a revise_quote_response
type ReviseQuoteResponse struct {
	ParentQuote *QuotingQuote `json:"parent_quote,omitempty"`
	Quote *QuotingQuote `json:"quote,omitempty"`
}
