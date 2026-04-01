package models

import (
	"time"
)

// Task represents a task
type Task struct {
	AssignedTo string `json:"assigned_to,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	Description string `json:"description,omitempty"`
	DueDate time.Time `json:"due_date,omitempty"`
	Id string `json:"id"`
	Priority interface{} `json:"priority"`
	RelatedEntityId string `json:"related_entity_id,omitempty"`
	RelatedEntityType string `json:"related_entity_type,omitempty"`
	Status interface{} `json:"status"`
	Title string `json:"title"`
	Type *TaskType `json:"type"`
}
