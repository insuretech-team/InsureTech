package models

import (
	"time"
)

// TrackPublishedEvent represents a track_published_event
type TrackPublishedEvent struct {
	RoomId string `json:"room_id,omitempty"`
	PeerId string `json:"peer_id,omitempty"`
	Track *Track `json:"track,omitempty"`
	PublishedAt time.Time `json:"published_at,omitempty"`
}
