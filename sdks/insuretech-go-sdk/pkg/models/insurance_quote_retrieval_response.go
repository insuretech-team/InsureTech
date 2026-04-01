package models


// InsuranceQuoteRetrievalResponse represents a insurance_quote_retrieval_response
type InsuranceQuoteRetrievalResponse struct {
	Quote *UnderwritingQuote `json:"quote,omitempty"`
}
