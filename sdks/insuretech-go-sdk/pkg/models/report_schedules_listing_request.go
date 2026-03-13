package models


// ReportSchedulesListingRequest represents a report_schedules_listing_request
type ReportSchedulesListingRequest struct {
	PageToken string `json:"page_token,omitempty"`
	ReportDefinitionId string `json:"report_definition_id"`
	ActiveOnly bool `json:"active_only,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}
