package models

import (
	"time"
)

// TrackPublishedEvent represents a track_published_event
type TrackPublishedEvent struct {
	PeerId string `json:"peer_id,omitempty"`
	PublishedAt time.Time `json:"published_at,omitempty"`
	RoomId string `json:"room_id,omitempty"`
	Track *Track `json:"track,omitempty"`
}
