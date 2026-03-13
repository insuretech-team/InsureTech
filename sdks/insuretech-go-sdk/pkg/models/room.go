package models

import (
	"time"
)

// Room represents a room
type Room struct {
	Name string `json:"name"`
	MaxParticipants int `json:"max_participants"`
	CreatedAt time.Time `json:"created_at"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	CreatorId string `json:"creator_id,omitempty"`
	RoomId string `json:"room_id"`
	Config interface{} `json:"config"`
	ParticipantCount int `json:"participant_count"`
	State interface{} `json:"state"`
	ClosedAt time.Time `json:"closed_at,omitempty"`
}
