package models

import (
	"time"
)

// RoomConfigUpdatedEvent represents a room_config_updated_event
type RoomConfigUpdatedEvent struct {
	RoomId string `json:"room_id,omitempty"`
	OldConfig *RoomConfig `json:"old_config,omitempty"`
	NewConfig *RoomConfig `json:"new_config,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
