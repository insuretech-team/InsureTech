package models


// LifeConvertQuoteToPolicyRequest represents a life_convert_quote_to_policy_request
type LifeConvertQuoteToPolicyRequest struct {
	ConvertedBy string `json:"converted_by,omitempty"`
	PolicyId string `json:"policy_id"`
	QuoteId string `json:"quote_id"`
}
