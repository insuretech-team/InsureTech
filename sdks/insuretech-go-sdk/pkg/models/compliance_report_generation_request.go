package models


// ComplianceReportGenerationRequest represents a compliance_report_generation_request
type ComplianceReportGenerationRequest struct {
	EndDate string `json:"end_date,omitempty"`
	Regulation string `json:"regulation"`
	StartDate string `json:"start_date,omitempty"`
}
