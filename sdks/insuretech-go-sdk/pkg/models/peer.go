package models

import (
	"time"
)

// Peer represents a peer
type Peer struct {
	DisplayName string `json:"display_name"`
	JoinedAt time.Time `json:"joined_at"`
	LastSeenAt time.Time `json:"last_seen_at"`
	LeftAt time.Time `json:"left_at,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	PeerId string `json:"peer_id"`
	RoomId string `json:"room_id"`
	State interface{} `json:"state"`
	Tracks []*Track `json:"tracks,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}
