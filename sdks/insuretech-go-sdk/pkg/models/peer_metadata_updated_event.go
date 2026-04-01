package models

import (
	"time"
)

// PeerMetadataUpdatedEvent represents a peer_metadata_updated_event
type PeerMetadataUpdatedEvent struct {
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	PeerId string `json:"peer_id,omitempty"`
	RoomId string `json:"room_id,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
