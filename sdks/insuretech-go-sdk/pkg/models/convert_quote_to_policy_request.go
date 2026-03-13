package models


// ConvertQuoteToPolicyRequest represents a convert_quote_to_policy_request
type ConvertQuoteToPolicyRequest struct {
	PaymentReference string `json:"payment_reference,omitempty"`
	QuoteId string `json:"quote_id"`
	PaymentMethod string `json:"payment_method,omitempty"`
}
