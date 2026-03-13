package models


// TaskCreationRequest represents a task_creation_request
type TaskCreationRequest struct {
	AssignedTo string `json:"assigned_to,omitempty"`
	RelatedEntityType string `json:"related_entity_type,omitempty"`
	RelatedEntityId string `json:"related_entity_id"`
	DueDate string `json:"due_date,omitempty"`
	Title string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Type string `json:"type"`
	Priority string `json:"priority,omitempty"`
}
