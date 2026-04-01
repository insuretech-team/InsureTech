package models

import (
	"time"
)

// PeerSession represents a peer_session
type PeerSession struct {
	JoinedAt time.Time `json:"joined_at"`
	LeftAt time.Time `json:"left_at,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	PeerId string `json:"peer_id"`
	PeerSessionId string `json:"peer_session_id"`
	SessionId string `json:"session_id"`
}
