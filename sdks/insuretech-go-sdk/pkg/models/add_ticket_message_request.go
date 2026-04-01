package models


// AddTicketMessageRequest represents a add_ticket_message_request
type AddTicketMessageRequest struct {
	Attachments string `json:"attachments,omitempty"`
	Content string `json:"content,omitempty"`
	IsInternal bool `json:"is_internal,omitempty"`
	TicketId string `json:"ticket_id"`
}
