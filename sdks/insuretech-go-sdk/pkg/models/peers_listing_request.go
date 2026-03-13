package models


// PeersListingRequest represents a peers_listing_request
type PeersListingRequest struct {
	RoomId string `json:"room_id"`
	StateFilter *PeerConnectionState `json:"state_filter,omitempty"`
}
