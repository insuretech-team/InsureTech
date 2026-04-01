package models

import (
	"time"
)

// RoomCreatedEvent represents a room_created_event
type RoomCreatedEvent struct {
	Config *RoomConfig `json:"config,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	CreatorId string `json:"creator_id,omitempty"`
	Name string `json:"name,omitempty"`
	RoomId string `json:"room_id,omitempty"`
}
