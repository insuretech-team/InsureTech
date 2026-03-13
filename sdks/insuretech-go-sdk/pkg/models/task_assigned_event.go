package models

import (
	"time"
)

// TaskAssignedEvent represents a task_assigned_event
type TaskAssignedEvent struct {
	EventId string `json:"event_id,omitempty"`
	TaskId string `json:"task_id,omitempty"`
	AssignedTo string `json:"assigned_to,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
