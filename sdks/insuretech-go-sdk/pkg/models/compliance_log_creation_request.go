package models


// ComplianceLogCreationRequest represents a compliance_log_creation_request
type ComplianceLogCreationRequest struct {
	Description string `json:"description,omitempty"`
	EntityId string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	Evidence string `json:"evidence,omitempty"`
	Regulation string `json:"regulation,omitempty"`
	Status string `json:"status,omitempty"`
	Type string `json:"type"`
}
