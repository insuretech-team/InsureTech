package models


// PaginationRequest represents a pagination_request
type PaginationRequest struct {
	Page int `json:"page"`
	PageSize int `json:"page_size,omitempty"`
	SortBy string `json:"sort_by,omitempty"`
	SortOrder *SortOrder `json:"sort_order,omitempty"`
}
