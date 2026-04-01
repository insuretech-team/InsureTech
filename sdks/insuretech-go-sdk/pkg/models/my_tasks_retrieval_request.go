package models


// MyTasksRetrievalRequest represents a my_tasks_retrieval_request
type MyTasksRetrievalRequest struct {
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	Status string `json:"status"`
}
