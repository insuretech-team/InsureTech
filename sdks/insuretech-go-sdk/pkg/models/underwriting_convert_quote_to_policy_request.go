package models


// UnderwritingConvertQuoteToPolicyRequest represents a underwriting_convert_quote_to_policy_request
type UnderwritingConvertQuoteToPolicyRequest struct {
	PaymentMethod string `json:"payment_method,omitempty"`
	PaymentReference string `json:"payment_reference,omitempty"`
	QuoteId string `json:"quote_id"`
}
