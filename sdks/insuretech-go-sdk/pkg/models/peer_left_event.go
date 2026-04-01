package models

import (
	"time"
)

// PeerLeftEvent represents a peer_left_event
type PeerLeftEvent struct {
	LeftAt time.Time `json:"left_at,omitempty"`
	PeerId string `json:"peer_id,omitempty"`
	Reason string `json:"reason,omitempty"`
	RoomId string `json:"room_id,omitempty"`
}
