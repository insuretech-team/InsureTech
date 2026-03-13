package models


// ComplianceReportGenerationRequest represents a compliance_report_generation_request
type ComplianceReportGenerationRequest struct {
	Regulation string `json:"regulation"`
	StartDate string `json:"start_date,omitempty"`
	EndDate string `json:"end_date,omitempty"`
}
