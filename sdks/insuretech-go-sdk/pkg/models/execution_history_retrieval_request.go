package models


// ExecutionHistoryRetrievalRequest represents a execution_history_retrieval_request
type ExecutionHistoryRetrievalRequest struct {
	BusinessWorkflowId string `json:"business_workflow_id"`
	EntityId string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	PageSize int `json:"page_size,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}
