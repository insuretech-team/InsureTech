package models


// BusinessWorkflowDeletionRequest represents a business_workflow_deletion_request
type BusinessWorkflowDeletionRequest struct {
	BusinessWorkflowId string `json:"business_workflow_id"`
	Permanent bool `json:"permanent,omitempty"`
}
