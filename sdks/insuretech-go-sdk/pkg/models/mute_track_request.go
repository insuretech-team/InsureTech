package models


// MuteTrackRequest represents a mute_track_request
type MuteTrackRequest struct {
	Muted bool `json:"muted,omitempty"`
	PeerId string `json:"peer_id"`
	RoomId string `json:"room_id"`
	TrackId string `json:"track_id"`
}
