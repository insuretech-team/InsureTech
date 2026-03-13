package models

import (
	"time"
)

// PeerMetadataUpdatedEvent represents a peer_metadata_updated_event
type PeerMetadataUpdatedEvent struct {
	RoomId string `json:"room_id,omitempty"`
	PeerId string `json:"peer_id,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
