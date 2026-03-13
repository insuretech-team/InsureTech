package models


// ReportExecutionsListingRequest represents a report_executions_listing_request
type ReportExecutionsListingRequest struct {
	EndDate string `json:"end_date,omitempty"`
	Page int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	ReportDefinitionId string `json:"report_definition_id"`
	Status string `json:"status,omitempty"`
	StartDate string `json:"start_date,omitempty"`
}
