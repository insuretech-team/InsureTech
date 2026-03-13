package models

import (
	"time"
)

// WorkflowCompletedEvent represents a workflow_completed_event
type WorkflowCompletedEvent struct {
	WorkflowInstanceId string `json:"workflow_instance_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	EntityId string `json:"entity_id,omitempty"`
	Status string `json:"status,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
}
