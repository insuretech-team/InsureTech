package models

import (
	"time"
)

// RoomClosedEvent represents a room_closed_event
type RoomClosedEvent struct {
	RoomId string `json:"room_id,omitempty"`
	Reason string `json:"reason,omitempty"`
	ClosedAt time.Time `json:"closed_at,omitempty"`
}
