package models


// MyTasksListingRequest represents a my_tasks_listing_request
type MyTasksListingRequest struct {
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	Priority string `json:"priority,omitempty"`
	Status string `json:"status"`
}
