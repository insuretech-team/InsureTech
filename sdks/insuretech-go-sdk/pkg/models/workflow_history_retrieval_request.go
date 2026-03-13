package models


// WorkflowHistoryRetrievalRequest represents a workflow_history_retrieval_request
type WorkflowHistoryRetrievalRequest struct {
	EntityType string `json:"entity_type"`
	EntityId string `json:"entity_id"`
}
