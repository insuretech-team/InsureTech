package models

// TaskTaskStatus represents a task_task_status
type TaskTaskStatus string

// TaskTaskStatus values
const (
	TaskTaskStatusTASKSTATUSUNSPECIFIED TaskTaskStatus = "TASK_STATUS_UNSPECIFIED"
	TaskTaskStatusTASKSTATUSPENDING  = "TASK_STATUS_PENDING"
	TaskTaskStatusTASKSTATUSINPROGRESS  = "TASK_STATUS_IN_PROGRESS"
	TaskTaskStatusTASKSTATUSCOMPLETED  = "TASK_STATUS_COMPLETED"
	TaskTaskStatusTASKSTATUSCANCELLED  = "TASK_STATUS_CANCELLED"
)
