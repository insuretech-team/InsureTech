package models

import (
	"time"
)

// PeerStateChangedEvent represents a peer_state_changed_event
type PeerStateChangedEvent struct {
	ChangedAt time.Time `json:"changed_at,omitempty"`
	NewState *PeerConnectionState `json:"new_state,omitempty"`
	OldState *PeerConnectionState `json:"old_state,omitempty"`
	PeerId string `json:"peer_id,omitempty"`
	RoomId string `json:"room_id,omitempty"`
}
