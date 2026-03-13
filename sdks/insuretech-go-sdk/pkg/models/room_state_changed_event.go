package models

import (
	"time"
)

// RoomStateChangedEvent represents a room_state_changed_event
type RoomStateChangedEvent struct {
	RoomId string `json:"room_id,omitempty"`
	OldState *RoomState `json:"old_state,omitempty"`
	NewState *RoomState `json:"new_state,omitempty"`
	ChangedAt time.Time `json:"changed_at,omitempty"`
}
