package models


// TicketStatusUpdateResponse represents a ticket_status_update_response
type TicketStatusUpdateResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
