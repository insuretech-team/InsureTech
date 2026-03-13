package models


// MuteTrackRequest represents a mute_track_request
type MuteTrackRequest struct {
	RoomId string `json:"room_id"`
	PeerId string `json:"peer_id"`
	TrackId string `json:"track_id"`
	Muted bool `json:"muted,omitempty"`
}
