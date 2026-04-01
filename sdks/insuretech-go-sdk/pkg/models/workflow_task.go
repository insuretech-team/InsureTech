package models

import (
	"time"
)

// WorkflowTask represents a workflow_task
type WorkflowTask struct {
	AssignedTo string `json:"assigned_to,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	Comments string `json:"comments,omitempty"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	Decision string `json:"decision,omitempty"`
	DueDate time.Time `json:"due_date,omitempty"`
	Id string `json:"id"`
	Status interface{} `json:"status"`
	StepName string `json:"step_name"`
	Type *WorkflowTaskType `json:"type"`
	WorkflowInstanceId string `json:"workflow_instance_id"`
}
