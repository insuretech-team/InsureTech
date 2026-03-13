package models


// RequestQuoteResponse represents a request_quote_response
type RequestQuoteResponse struct {
	TotalPremium *Money `json:"total_premium,omitempty"`
	ValidUntil string `json:"valid_until,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
	QuoteNumber string `json:"quote_number,omitempty"`
	BasePremium *Money `json:"base_premium,omitempty"`
}
