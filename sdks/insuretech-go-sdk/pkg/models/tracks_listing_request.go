package models


// TracksListingRequest represents a tracks_listing_request
type TracksListingRequest struct {
	PeerId string `json:"peer_id"`
	RoomId string `json:"room_id"`
	TypeFilter *TrackType `json:"type_filter,omitempty"`
}
