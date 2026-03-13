package models

import (
	"time"
)

// RenegotiationRequiredEvent represents a renegotiation_required_event
type RenegotiationRequiredEvent struct {
	RoomId string `json:"room_id,omitempty"`
	PeerId string `json:"peer_id,omitempty"`
	Reason string `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
