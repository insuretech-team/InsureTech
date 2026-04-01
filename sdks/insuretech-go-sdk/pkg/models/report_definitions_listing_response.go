package models


// ReportDefinitionsListingResponse represents a report_definitions_listing_response
type ReportDefinitionsListingResponse struct {
	NextPageToken string `json:"next_page_token,omitempty"`
	ReportDefinitions []*ReportDefinition `json:"report_definitions,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
