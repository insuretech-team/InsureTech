package models

import (
	"time"
)

// TaskCreatedEvent represents a task_created_event
type TaskCreatedEvent struct {
	AssignedTo string `json:"assigned_to,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	TaskId string `json:"task_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Title string `json:"title,omitempty"`
	Type string `json:"type,omitempty"`
}
