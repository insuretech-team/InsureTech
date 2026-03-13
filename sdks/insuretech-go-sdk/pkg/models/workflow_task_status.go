package models

// WorkflowTaskStatus represents a workflow_task_status
type WorkflowTaskStatus string

// WorkflowTaskStatus values
const (
	WorkflowTaskStatusTASKSTATUSUNSPECIFIED WorkflowTaskStatus = "TASK_STATUS_UNSPECIFIED"
	WorkflowTaskStatusTASKSTATUSPENDING  = "TASK_STATUS_PENDING"
	WorkflowTaskStatusTASKSTATUSINPROGRESS  = "TASK_STATUS_IN_PROGRESS"
	WorkflowTaskStatusTASKSTATUSCOMPLETED  = "TASK_STATUS_COMPLETED"
	WorkflowTaskStatusTASKSTATUSSKIPPED  = "TASK_STATUS_SKIPPED"
)
