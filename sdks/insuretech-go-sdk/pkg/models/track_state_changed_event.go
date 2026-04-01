package models

import (
	"time"
)

// TrackStateChangedEvent represents a track_state_changed_event
type TrackStateChangedEvent struct {
	ChangedAt time.Time `json:"changed_at,omitempty"`
	NewState *TrackState `json:"new_state,omitempty"`
	OldState *TrackState `json:"old_state,omitempty"`
	PeerId string `json:"peer_id,omitempty"`
	RoomId string `json:"room_id,omitempty"`
	TrackId string `json:"track_id,omitempty"`
}
