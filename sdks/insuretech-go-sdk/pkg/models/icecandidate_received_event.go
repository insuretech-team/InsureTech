package models

import (
	"time"
)

// ICECandidateReceivedEvent represents a icecandidate_received_event
type ICECandidateReceivedEvent struct {
	RoomId string `json:"room_id,omitempty"`
	Candidate *ICECandidate `json:"candidate,omitempty"`
	ReceivedAt time.Time `json:"received_at,omitempty"`
}
