package models


// PageRequest represents a page_request
type PageRequest struct {
	Page int `json:"page"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}
