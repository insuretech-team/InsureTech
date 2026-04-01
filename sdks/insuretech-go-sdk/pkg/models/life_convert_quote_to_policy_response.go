package models

import (
	"time"
)

// LifeConvertQuoteToPolicyResponse represents a life_convert_quote_to_policy_response
type LifeConvertQuoteToPolicyResponse struct {
	ConvertedAt time.Time `json:"converted_at,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
}
