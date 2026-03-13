package models


// QuoteRetrievalResponse represents a quote_retrieval_response
type QuoteRetrievalResponse struct {
	Quote *Quote `json:"quote,omitempty"`
	HealthDeclaration *UnderwritingHealthDeclaration `json:"health_declaration,omitempty"`
	Decision *UnderwritingDecision `json:"decision,omitempty"`
	Error *Error `json:"error,omitempty"`
}
