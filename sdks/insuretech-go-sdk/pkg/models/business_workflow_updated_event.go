package models

import (
	"time"
)

// BusinessWorkflowUpdatedEvent represents a business_workflow_updated_event
type BusinessWorkflowUpdatedEvent struct {
	BusinessWorkflowId string `json:"business_workflow_id,omitempty"`
	ChangedFields []string `json:"changed_fields,omitempty"`
	EventId string `json:"event_id,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
	Version int `json:"version,omitempty"`
	WorkflowName string `json:"workflow_name,omitempty"`
}
