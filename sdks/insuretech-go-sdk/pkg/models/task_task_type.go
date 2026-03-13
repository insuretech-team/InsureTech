package models

// TaskTaskType represents a task_task_type
type TaskTaskType string

// TaskTaskType values
const (
	TaskTaskTypeTASKTYPEUNSPECIFIED TaskTaskType = "TASK_TYPE_UNSPECIFIED"
	TaskTaskTypeTASKTYPEFOLLOWUP  = "TASK_TYPE_FOLLOW_UP"
	TaskTaskTypeTASKTYPEREVIEW  = "TASK_TYPE_REVIEW"
	TaskTaskTypeTASKTYPEAPPROVAL  = "TASK_TYPE_APPROVAL"
	TaskTaskTypeTASKTYPEDOCUMENT  = "TASK_TYPE_DOCUMENT"
	TaskTaskTypeTASKTYPEREMINDER  = "TASK_TYPE_REMINDER"
)
