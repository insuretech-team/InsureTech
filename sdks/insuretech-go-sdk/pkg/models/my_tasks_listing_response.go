package models


// MyTasksListingResponse represents a my_tasks_listing_response
type MyTasksListingResponse struct {
	Tasks []*Task `json:"tasks,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
