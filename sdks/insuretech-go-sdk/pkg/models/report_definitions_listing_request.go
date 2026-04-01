package models


// ReportDefinitionsListingRequest represents a report_definitions_listing_request
type ReportDefinitionsListingRequest struct {
	ActiveOnly bool `json:"active_only,omitempty"`
	Category string `json:"category"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}
