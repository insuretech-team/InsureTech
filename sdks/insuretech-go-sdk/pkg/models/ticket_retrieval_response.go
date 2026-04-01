package models


// TicketRetrievalResponse represents a ticket_retrieval_response
type TicketRetrievalResponse struct {
	Messages []*TicketMessage `json:"messages,omitempty"`
	Ticket *Ticket `json:"ticket,omitempty"`
}
