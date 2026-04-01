package models


// TicketMessage represents a ticket_message
type TicketMessage struct {
	Attachments string `json:"attachments,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	Content string `json:"content"`
	Id string `json:"id"`
	IsInternal bool `json:"is_internal,omitempty"`
	SenderId string `json:"sender_id"`
	TicketId string `json:"ticket_id"`
	Type *MessageType `json:"type"`
}
