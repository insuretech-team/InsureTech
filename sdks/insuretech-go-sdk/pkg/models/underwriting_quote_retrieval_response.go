package models


// UnderwritingQuoteRetrievalResponse represents a underwriting_quote_retrieval_response
type UnderwritingQuoteRetrievalResponse struct {
	Quote *Quote `json:"quote,omitempty"`
	HealthDeclaration *UnderwritingHealthDeclaration `json:"health_declaration,omitempty"`
	Decision *UnderwritingDecision `json:"decision,omitempty"`
	Error *Error `json:"error,omitempty"`
}
