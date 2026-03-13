package models


// TicketRetrievalResponse represents a ticket_retrieval_response
type TicketRetrievalResponse struct {
	Ticket *Ticket `json:"ticket,omitempty"`
	Messages []*TicketMessage `json:"messages,omitempty"`
	Error *Error `json:"error,omitempty"`
}
