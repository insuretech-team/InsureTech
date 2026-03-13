package models


// TracksListingRequest represents a tracks_listing_request
type TracksListingRequest struct {
	RoomId string `json:"room_id"`
	PeerId string `json:"peer_id"`
	TypeFilter *TrackType `json:"type_filter,omitempty"`
}
