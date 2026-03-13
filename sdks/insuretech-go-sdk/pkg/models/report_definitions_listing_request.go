package models


// ReportDefinitionsListingRequest represents a report_definitions_listing_request
type ReportDefinitionsListingRequest struct {
	Category string `json:"category"`
	ActiveOnly bool `json:"active_only,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}
