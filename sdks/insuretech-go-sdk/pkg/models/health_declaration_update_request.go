package models


// HealthDeclarationUpdateRequest represents a health_declaration_update_request
type HealthDeclarationUpdateRequest struct {
	Declaration *UnderwritingHealthDeclaration `json:"declaration"`
}
