package models


// HealthDeclarationCreationRequest represents a health_declaration_creation_request
type HealthDeclarationCreationRequest struct {
	Declaration *UnderwritingHealthDeclaration `json:"declaration"`
}
