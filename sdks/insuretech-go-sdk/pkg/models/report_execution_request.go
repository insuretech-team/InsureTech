package models


// ReportExecutionRequest represents a report_execution_request
type ReportExecutionRequest struct {
	Format string `json:"format,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	ReportDefinitionId string `json:"report_definition_id"`
}
