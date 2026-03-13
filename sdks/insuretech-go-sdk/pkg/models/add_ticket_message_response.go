package models


// AddTicketMessageResponse represents a add_ticket_message_response
type AddTicketMessageResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
	MessageId string `json:"message_id,omitempty"`
}
