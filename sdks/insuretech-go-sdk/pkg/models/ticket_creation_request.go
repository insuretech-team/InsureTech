package models


// TicketCreationRequest represents a ticket_creation_request
type TicketCreationRequest struct {
	BeneficiaryId string `json:"beneficiary_id"`
	Category string `json:"category,omitempty"`
	Description string `json:"description,omitempty"`
	Priority string `json:"priority,omitempty"`
	RelatedEntityId string `json:"related_entity_id"`
	RelatedEntityType string `json:"related_entity_type,omitempty"`
	Subject string `json:"subject,omitempty"`
	Type string `json:"type"`
}
