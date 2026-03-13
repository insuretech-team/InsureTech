package models


// HealthDeclarationRetrievalResponse represents a health_declaration_retrieval_response
type HealthDeclarationRetrievalResponse struct {
	HealthDeclaration *UnderwritingHealthDeclaration `json:"health_declaration,omitempty"`
	Error *Error `json:"error,omitempty"`
}
