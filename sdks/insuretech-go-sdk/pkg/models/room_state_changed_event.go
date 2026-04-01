package models

import (
	"time"
)

// RoomStateChangedEvent represents a room_state_changed_event
type RoomStateChangedEvent struct {
	ChangedAt time.Time `json:"changed_at,omitempty"`
	NewState *RoomState `json:"new_state,omitempty"`
	OldState *RoomState `json:"old_state,omitempty"`
	RoomId string `json:"room_id,omitempty"`
}
