package models


// HealthDeclarationCreationResponse represents a health_declaration_creation_response
type HealthDeclarationCreationResponse struct {
	Declaration *UnderwritingHealthDeclaration `json:"declaration,omitempty"`
}
