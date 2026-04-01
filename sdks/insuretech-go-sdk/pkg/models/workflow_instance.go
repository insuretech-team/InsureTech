package models

import (
	"time"
)

// WorkflowInstance represents a workflow_instance
type WorkflowInstance struct {
	AuditInfo interface{} `json:"audit_info"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	Context string `json:"context,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	CurrentStep string `json:"current_step,omitempty"`
	EntityId string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	Id string `json:"id"`
	InitiatedBy string `json:"initiated_by"`
	StartedAt time.Time `json:"started_at"`
	Status interface{} `json:"status"`
	WorkflowDefinitionId string `json:"workflow_definition_id"`
}
