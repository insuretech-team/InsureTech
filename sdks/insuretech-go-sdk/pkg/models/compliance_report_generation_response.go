package models


// ComplianceReportGenerationResponse represents a compliance_report_generation_response
type ComplianceReportGenerationResponse struct {
	ReportUrl string `json:"report_url,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
