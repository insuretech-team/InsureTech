package models

import (
	"time"
)

// RoomCreatedEvent represents a room_created_event
type RoomCreatedEvent struct {
	RoomId string `json:"room_id,omitempty"`
	Name string `json:"name,omitempty"`
	CreatorId string `json:"creator_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Config *RoomConfig `json:"config,omitempty"`
}
