package models

import (
	"time"
)

// BusinessWorkflowDeletedEvent represents a business_workflow_deleted_event
type BusinessWorkflowDeletedEvent struct {
	BusinessWorkflowId string `json:"business_workflow_id,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	DeletedBy string `json:"deleted_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Permanent bool `json:"permanent,omitempty"`
	WorkflowName string `json:"workflow_name,omitempty"`
}
