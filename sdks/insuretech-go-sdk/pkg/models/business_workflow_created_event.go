package models

import (
	"time"
)

// BusinessWorkflowCreatedEvent represents a business_workflow_created_event
type BusinessWorkflowCreatedEvent struct {
	BusinessWorkflowId string `json:"business_workflow_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	WorkflowName string `json:"workflow_name,omitempty"`
	WorkflowType *BusinessWorkflowType `json:"workflow_type,omitempty"`
}
