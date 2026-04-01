package models

import (
	"time"
)

// RoomSession represents a room_session
type RoomSession struct {
	EndedAt time.Time `json:"ended_at,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	ParticipantCount int `json:"participant_count"`
	RoomId string `json:"room_id"`
	SessionId string `json:"session_id"`
	StartedAt time.Time `json:"started_at"`
}
