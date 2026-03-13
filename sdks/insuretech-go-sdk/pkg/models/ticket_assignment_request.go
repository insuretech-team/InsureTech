package models


// TicketAssignmentRequest represents a ticket_assignment_request
type TicketAssignmentRequest struct {
	TicketId string `json:"ticket_id"`
	AssignedTo string `json:"assigned_to,omitempty"`
}
