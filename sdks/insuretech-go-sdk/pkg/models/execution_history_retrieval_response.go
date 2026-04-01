package models


// ExecutionHistoryRetrievalResponse represents a execution_history_retrieval_response
type ExecutionHistoryRetrievalResponse struct {
	Executions []*BusinessWorkflowExecution `json:"executions,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
