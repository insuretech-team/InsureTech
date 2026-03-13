package models


// PublishTrackRequest represents a publish_track_request
type PublishTrackRequest struct {
	RoomId string `json:"room_id"`
	PeerId string `json:"peer_id"`
	Track *Track `json:"track,omitempty"`
}
