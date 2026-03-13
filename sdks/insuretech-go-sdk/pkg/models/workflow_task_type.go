package models

// WorkflowTaskType represents a workflow_task_type
type WorkflowTaskType string

// WorkflowTaskType values
const (
	WorkflowTaskTypeTASKTYPEUNSPECIFIED WorkflowTaskType = "TASK_TYPE_UNSPECIFIED"
	WorkflowTaskTypeTASKTYPEAPPROVAL  = "TASK_TYPE_APPROVAL"
	WorkflowTaskTypeTASKTYPEREVIEW  = "TASK_TYPE_REVIEW"
	WorkflowTaskTypeTASKTYPENOTIFICATION  = "TASK_TYPE_NOTIFICATION"
	WorkflowTaskTypeTASKTYPEACTION  = "TASK_TYPE_ACTION"
)
