package models


// PeersListingResponse represents a peers_listing_response
type PeersListingResponse struct {
	Peers []*Peer `json:"peers,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}
