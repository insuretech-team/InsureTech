package models


// PaginationResponse represents a pagination_response
type PaginationResponse struct {
	HasNext bool `json:"has_next,omitempty"`
	HasPrevious bool `json:"has_previous,omitempty"`
	TotalItems int `json:"total_items,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
	CurrentPage int `json:"current_page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}
