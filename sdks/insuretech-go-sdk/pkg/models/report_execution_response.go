package models


// ReportExecutionResponse represents a report_execution_response
type ReportExecutionResponse struct {
	ReportExecutionId string `json:"report_execution_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
