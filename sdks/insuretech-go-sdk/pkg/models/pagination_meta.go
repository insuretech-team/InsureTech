package models


// PaginationMeta represents a pagination_meta
type PaginationMeta struct {
	HasNext bool `json:"has_next"`
	HasPrevious bool `json:"has_previous"`
	NextPageToken string `json:"next_page_token,omitempty"`
	Page int `json:"page"`
	PageSize int `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int `json:"total_pages"`
}
