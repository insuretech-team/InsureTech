package models

import (
	"time"
)

// PeerLeftEvent represents a peer_left_event
type PeerLeftEvent struct {
	RoomId string `json:"room_id,omitempty"`
	PeerId string `json:"peer_id,omitempty"`
	Reason string `json:"reason,omitempty"`
	LeftAt time.Time `json:"left_at,omitempty"`
}
