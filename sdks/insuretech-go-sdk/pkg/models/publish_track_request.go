package models


// PublishTrackRequest represents a publish_track_request
type PublishTrackRequest struct {
	PeerId string `json:"peer_id"`
	RoomId string `json:"room_id"`
	Track *Track `json:"track,omitempty"`
}
