package models


// TicketsListingResponse represents a tickets_listing_response
type TicketsListingResponse struct {
	Tickets []*Ticket `json:"tickets,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	Error *Error `json:"error,omitempty"`
}
