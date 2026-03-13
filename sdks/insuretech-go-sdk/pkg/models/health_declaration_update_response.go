package models


// HealthDeclarationUpdateResponse represents a health_declaration_update_response
type HealthDeclarationUpdateResponse struct {
	Declaration *UnderwritingHealthDeclaration `json:"declaration,omitempty"`
}
