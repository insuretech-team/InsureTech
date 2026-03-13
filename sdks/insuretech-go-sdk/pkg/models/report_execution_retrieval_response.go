package models


// ReportExecutionRetrievalResponse represents a report_execution_retrieval_response
type ReportExecutionRetrievalResponse struct {
	ReportExecution *ReportExecution `json:"report_execution,omitempty"`
	Error *Error `json:"error,omitempty"`
}
