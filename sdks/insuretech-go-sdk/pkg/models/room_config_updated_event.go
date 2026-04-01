package models

import (
	"time"
)

// RoomConfigUpdatedEvent represents a room_config_updated_event
type RoomConfigUpdatedEvent struct {
	NewConfig *RoomConfig `json:"new_config,omitempty"`
	OldConfig *RoomConfig `json:"old_config,omitempty"`
	RoomId string `json:"room_id,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
