package models


// RequestQuoteResponse represents a request_quote_response
type RequestQuoteResponse struct {
	BasePremium *Money `json:"base_premium,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	QuoteNumber string `json:"quote_number,omitempty"`
	TotalPremium *Money `json:"total_premium,omitempty"`
	ValidUntil string `json:"valid_until,omitempty"`
}
