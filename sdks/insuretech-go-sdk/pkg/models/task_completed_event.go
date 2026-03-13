package models

import (
	"time"
)

// TaskCompletedEvent represents a task_completed_event
type TaskCompletedEvent struct {
	EventId string `json:"event_id,omitempty"`
	TaskId string `json:"task_id,omitempty"`
	CompletedBy string `json:"completed_by,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
