package models


// UnderwritingHealthDeclarationRetrievalResponse represents a underwriting_health_declaration_retrieval_response
type UnderwritingHealthDeclarationRetrievalResponse struct {
	Error *Error `json:"error,omitempty"`
	HealthDeclaration *UnderwritingHealthDeclaration `json:"health_declaration,omitempty"`
}
