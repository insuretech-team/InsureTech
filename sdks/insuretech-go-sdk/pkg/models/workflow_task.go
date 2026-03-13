package models

import (
	"time"
)

// WorkflowTask represents a workflow_task
type WorkflowTask struct {
	WorkflowInstanceId string `json:"workflow_instance_id"`
	Type *WorkflowTaskType `json:"type"`
	AssignedTo string `json:"assigned_to,omitempty"`
	Status interface{} `json:"status"`
	DueDate time.Time `json:"due_date,omitempty"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	Id string `json:"id"`
	StepName string `json:"step_name"`
	Decision string `json:"decision,omitempty"`
	Comments string `json:"comments,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
}
