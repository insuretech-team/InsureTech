package models


// MyTasksListingRequest represents a my_tasks_listing_request
type MyTasksListingRequest struct {
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	Status string `json:"status"`
	Priority string `json:"priority,omitempty"`
}
