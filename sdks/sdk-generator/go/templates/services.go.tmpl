package services

import (
	"context"
	"fmt"
)

// Client defines the interface for making HTTP requests
type Client interface {
	DoRequest(ctx context.Context, method, path string, body, result interface{}) error
}

// ListOptions represents common list query parameters
type ListOptions struct {
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
	SortBy   string `json:"sort_by,omitempty"`
	SortDir  string `json:"sort_dir,omitempty"`
}

// ListResponse represents a paginated list response
type ListResponse struct {
	Items      []interface{} `json:"items"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalItems int           `json:"total_items"`
	TotalPages int           `json:"total_pages"`
	HasMore    bool          `json:"has_more"`
}

// buildListPath builds a URL path with list query parameters
func buildListPath(basePath string, opts *ListOptions) string {
	if opts == nil {
		return basePath
	}
	
	path := basePath + "?"
	if opts.Page > 0 {
		path += fmt.Sprintf("page=%d&", opts.Page)
	}
	if opts.PageSize > 0 {
		path += fmt.Sprintf("page_size=%d&", opts.PageSize)
	}
	if opts.SortBy != "" {
		path += fmt.Sprintf("sort_by=%s&", opts.SortBy)
	}
	if opts.SortDir != "" {
		path += fmt.Sprintf("sort_dir=%s&", opts.SortDir)
	}
	
	return path
}
