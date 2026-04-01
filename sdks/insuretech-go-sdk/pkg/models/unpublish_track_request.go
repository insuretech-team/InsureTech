package models


// UnpublishTrackRequest represents a unpublish_track_request
type UnpublishTrackRequest struct {
	PeerId string `json:"peer_id"`
	RoomId string `json:"room_id"`
	TrackId string `json:"track_id"`
}
