package models

import (
	"time"
)

// RoomSession represents a room_session
type RoomSession struct {
	SessionId string `json:"session_id"`
	RoomId string `json:"room_id"`
	StartedAt time.Time `json:"started_at"`
	EndedAt time.Time `json:"ended_at,omitempty"`
	ParticipantCount int `json:"participant_count"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}
