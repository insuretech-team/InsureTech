package models


// ReportGenerationRequest represents a report_generation_request
type ReportGenerationRequest struct {
	EndDate string `json:"end_date,omitempty"`
	ReportId string `json:"report_id"`
	StartDate string `json:"start_date,omitempty"`
}
