package models


// WorkflowHistoryRetrievalRequest represents a workflow_history_retrieval_request
type WorkflowHistoryRetrievalRequest struct {
	EntityId string `json:"entity_id"`
	EntityType string `json:"entity_type"`
}
