package models


// TicketAssignmentRequest represents a ticket_assignment_request
type TicketAssignmentRequest struct {
	AssignedTo string `json:"assigned_to,omitempty"`
	TicketId string `json:"ticket_id"`
}
