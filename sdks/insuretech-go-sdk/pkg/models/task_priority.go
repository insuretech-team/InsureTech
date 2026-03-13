package models

// TaskPriority represents a task_priority
type TaskPriority string

// TaskPriority values
const (
	TaskPriorityTASKPRIORITYUNSPECIFIED TaskPriority = "TASK_PRIORITY_UNSPECIFIED"
	TaskPriorityTASKPRIORITYLOW  = "TASK_PRIORITY_LOW"
	TaskPriorityTASKPRIORITYMEDIUM  = "TASK_PRIORITY_MEDIUM"
	TaskPriorityTASKPRIORITYHIGH  = "TASK_PRIORITY_HIGH"
	TaskPriorityTASKPRIORITYURGENT  = "TASK_PRIORITY_URGENT"
)
