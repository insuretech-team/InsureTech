package models


// HealthDeclarationByQuoteRetrievalResponse represents a health_declaration_by_quote_retrieval_response
type HealthDeclarationByQuoteRetrievalResponse struct {
	Declaration *UnderwritingHealthDeclaration `json:"declaration,omitempty"`
}
