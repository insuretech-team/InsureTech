package models


// TicketCreationResponse represents a ticket_creation_response
type TicketCreationResponse struct {
	Error *Error `json:"error,omitempty"`
	TicketId string `json:"ticket_id,omitempty"`
	TicketNumber string `json:"ticket_number,omitempty"`
	Message string `json:"message,omitempty"`
}
