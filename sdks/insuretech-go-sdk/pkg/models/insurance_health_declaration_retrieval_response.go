package models


// InsuranceHealthDeclarationRetrievalResponse represents a insurance_health_declaration_retrieval_response
type InsuranceHealthDeclarationRetrievalResponse struct {
	Declaration *UnderwritingHealthDeclaration `json:"declaration,omitempty"`
}
