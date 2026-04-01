package models


// TaskUpdateRequest represents a task_update_request
type TaskUpdateRequest struct {
	Description string `json:"description,omitempty"`
	DueDate string `json:"due_date,omitempty"`
	Priority string `json:"priority,omitempty"`
	Status string `json:"status,omitempty"`
	TaskId string `json:"task_id"`
	Title string `json:"title,omitempty"`
}
