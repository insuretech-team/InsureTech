package models


// TaskRetrievalResponse represents a task_retrieval_response
type TaskRetrievalResponse struct {
	Task *Task `json:"task,omitempty"`
	Error *Error `json:"error,omitempty"`
}
