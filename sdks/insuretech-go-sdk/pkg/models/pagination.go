package models

// PaginatedResult represents a generic paginated response
type PaginatedResult[T any] struct {
	Items      []T   `json:"items"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int   `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	HasMore    bool  `json:"has_more"`
}

// PaginationOptions represents common pagination parameters
type PaginationOptions struct {
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
	SortBy   string `json:"sort_by,omitempty"`
	SortDir  string `json:"sort_dir,omitempty"`
}

// DefaultPaginationOptions returns default pagination options
func DefaultPaginationOptions() *PaginationOptions {
	return &PaginationOptions{
		Page:     1,
		PageSize: 20,
		SortDir:  "asc",
	}
}
