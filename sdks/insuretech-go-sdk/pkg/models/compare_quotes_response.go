package models


// CompareQuotesResponse represents a compare_quotes_response
type CompareQuotesResponse struct {
	Comparisons []*QuoteComparison `json:"comparisons,omitempty"`
}
