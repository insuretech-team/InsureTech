package models

import (
	"time"
)

// PeerJoinedEvent represents a peer_joined_event
type PeerJoinedEvent struct {
	Peer *Peer `json:"peer,omitempty"`
	JoinedAt time.Time `json:"joined_at,omitempty"`
	RoomId string `json:"room_id,omitempty"`
}
