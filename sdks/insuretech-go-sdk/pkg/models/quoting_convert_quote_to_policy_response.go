package models

import (
	"time"
)

// QuotingConvertQuoteToPolicyResponse represents a quoting_convert_quote_to_policy_response
type QuotingConvertQuoteToPolicyResponse struct {
	ConvertedAt time.Time `json:"converted_at,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	QuoteId string `json:"quote_id,omitempty"`
}
