package models


// MyTasksRetrievalRequest represents a my_tasks_retrieval_request
type MyTasksRetrievalRequest struct {
	Status string `json:"status"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}
