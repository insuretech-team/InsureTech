package models

import (
	"time"
)

// WorkflowCompletedEvent represents a workflow_completed_event
type WorkflowCompletedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	EntityId string `json:"entity_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Status string `json:"status,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	WorkflowInstanceId string `json:"workflow_instance_id,omitempty"`
}
