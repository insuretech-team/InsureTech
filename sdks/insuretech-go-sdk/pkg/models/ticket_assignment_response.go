package models


// TicketAssignmentResponse represents a ticket_assignment_response
type TicketAssignmentResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
