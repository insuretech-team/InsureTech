package models


// TicketCreationRequest represents a ticket_creation_request
type TicketCreationRequest struct {
	Category string `json:"category,omitempty"`
	Priority string `json:"priority,omitempty"`
	Subject string `json:"subject,omitempty"`
	Description string `json:"description,omitempty"`
	RelatedEntityType string `json:"related_entity_type,omitempty"`
	RelatedEntityId string `json:"related_entity_id"`
	BeneficiaryId string `json:"beneficiary_id"`
	Type string `json:"type"`
}
