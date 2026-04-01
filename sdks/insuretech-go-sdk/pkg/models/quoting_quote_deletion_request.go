package models


// QuotingQuoteDeletionRequest represents a quoting_quote_deletion_request
type QuotingQuoteDeletionRequest struct {
	Permanent bool `json:"permanent,omitempty"`
	QuoteId string `json:"quote_id"`
}
