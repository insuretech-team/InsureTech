package models


// PageResponse represents a page_response
type PageResponse struct {
	NextPageToken string `json:"next_page_token,omitempty"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	TotalItems int64 `json:"total_items,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}
