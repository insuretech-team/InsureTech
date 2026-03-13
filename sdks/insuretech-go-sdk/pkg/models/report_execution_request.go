package models


// ReportExecutionRequest represents a report_execution_request
type ReportExecutionRequest struct {
	ReportDefinitionId string `json:"report_definition_id"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	Format string `json:"format,omitempty"`
}
