package models


// TicketStatusUpdateRequest represents a ticket_status_update_request
type TicketStatusUpdateRequest struct {
	Comments string `json:"comments,omitempty"`
	Status string `json:"status,omitempty"`
	TicketId string `json:"ticket_id"`
}
