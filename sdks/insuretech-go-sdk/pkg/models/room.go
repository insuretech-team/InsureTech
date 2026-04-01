package models

import (
	"time"
)

// Room represents a room
type Room struct {
	ClosedAt time.Time `json:"closed_at,omitempty"`
	Config interface{} `json:"config"`
	CreatedAt time.Time `json:"created_at"`
	CreatorId string `json:"creator_id,omitempty"`
	MaxParticipants int `json:"max_participants"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Name string `json:"name"`
	ParticipantCount int `json:"participant_count"`
	RoomId string `json:"room_id"`
	State interface{} `json:"state"`
}
