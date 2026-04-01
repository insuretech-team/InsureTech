package models


// UnderwritingQuoteRetrievalResponse represents a underwriting_quote_retrieval_response
type UnderwritingQuoteRetrievalResponse struct {
	Decision *UnderwritingDecision `json:"decision,omitempty"`
	HealthDeclaration *UnderwritingHealthDeclaration `json:"health_declaration,omitempty"`
	Quote *UnderwritingQuote `json:"quote,omitempty"`
}
