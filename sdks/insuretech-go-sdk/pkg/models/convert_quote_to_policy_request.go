package models


// ConvertQuoteToPolicyRequest represents a convert_quote_to_policy_request
type ConvertQuoteToPolicyRequest struct {
	PaymentMethod string `json:"payment_method,omitempty"`
	PaymentReference string `json:"payment_reference,omitempty"`
	QuoteId string `json:"quote_id"`
}
