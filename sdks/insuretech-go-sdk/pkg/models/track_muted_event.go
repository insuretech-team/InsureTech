package models

import (
	"time"
)

// TrackMutedEvent represents a track_muted_event
type TrackMutedEvent struct {
	ChangedAt time.Time `json:"changed_at,omitempty"`
	RoomId string `json:"room_id,omitempty"`
	PeerId string `json:"peer_id,omitempty"`
	TrackId string `json:"track_id,omitempty"`
	Muted bool `json:"muted,omitempty"`
}
