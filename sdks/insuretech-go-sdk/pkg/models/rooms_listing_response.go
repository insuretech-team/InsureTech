package models


// RoomsListingResponse represents a rooms_listing_response
type RoomsListingResponse struct {
	Rooms []*Room `json:"rooms,omitempty"`
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
