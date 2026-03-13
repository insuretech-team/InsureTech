package models

// TicketPriority represents a ticket_priority
type TicketPriority string

// TicketPriority values
const (
	TicketPriorityTICKETPRIORITYUNSPECIFIED TicketPriority = "TICKET_PRIORITY_UNSPECIFIED"
	TicketPriorityTICKETPRIORITYLOW  = "TICKET_PRIORITY_LOW"
	TicketPriorityTICKETPRIORITYMEDIUM  = "TICKET_PRIORITY_MEDIUM"
	TicketPriorityTICKETPRIORITYHIGH  = "TICKET_PRIORITY_HIGH"
	TicketPriorityTICKETPRIORITYURGENT  = "TICKET_PRIORITY_URGENT"
)
