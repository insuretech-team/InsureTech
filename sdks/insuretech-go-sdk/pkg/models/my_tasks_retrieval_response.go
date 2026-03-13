package models


// MyTasksRetrievalResponse represents a my_tasks_retrieval_response
type MyTasksRetrievalResponse struct {
	Tasks []*WorkflowTask `json:"tasks,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
