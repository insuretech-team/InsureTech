package models


// AddTicketMessageRequest represents a add_ticket_message_request
type AddTicketMessageRequest struct {
	Content string `json:"content,omitempty"`
	IsInternal bool `json:"is_internal,omitempty"`
	Attachments string `json:"attachments,omitempty"`
	TicketId string `json:"ticket_id"`
}
