package models

import (
	"time"
)

// WorkflowTaskCompletedEvent represents a workflow_task_completed_event
type WorkflowTaskCompletedEvent struct {
	EventId string `json:"event_id,omitempty"`
	TaskId string `json:"task_id,omitempty"`
	WorkflowInstanceId string `json:"workflow_instance_id,omitempty"`
	Decision string `json:"decision,omitempty"`
	CompletedBy string `json:"completed_by,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
