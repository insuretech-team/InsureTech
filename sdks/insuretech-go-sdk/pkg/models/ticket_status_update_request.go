package models


// TicketStatusUpdateRequest represents a ticket_status_update_request
type TicketStatusUpdateRequest struct {
	TicketId string `json:"ticket_id"`
	Status string `json:"status,omitempty"`
	Comments string `json:"comments,omitempty"`
}
