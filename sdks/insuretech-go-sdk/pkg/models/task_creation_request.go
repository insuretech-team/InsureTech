package models


// TaskCreationRequest represents a task_creation_request
type TaskCreationRequest struct {
	AssignedTo string `json:"assigned_to,omitempty"`
	Description string `json:"description,omitempty"`
	DueDate string `json:"due_date,omitempty"`
	Priority string `json:"priority,omitempty"`
	RelatedEntityId string `json:"related_entity_id"`
	RelatedEntityType string `json:"related_entity_type,omitempty"`
	Title string `json:"title,omitempty"`
	Type string `json:"type"`
}
