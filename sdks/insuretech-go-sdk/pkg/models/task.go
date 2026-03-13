package models

import (
	"time"
)

// Task represents a task
type Task struct {
	Priority interface{} `json:"priority"`
	CreatedBy string `json:"created_by,omitempty"`
	RelatedEntityType string `json:"related_entity_type,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	Id string `json:"id"`
	Title string `json:"title"`
	Type *TaskTaskType `json:"type"`
	Status interface{} `json:"status"`
	AssignedTo string `json:"assigned_to,omitempty"`
	RelatedEntityId string `json:"related_entity_id,omitempty"`
	DueDate time.Time `json:"due_date,omitempty"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	Description string `json:"description,omitempty"`
}
