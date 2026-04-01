package models

import (
	"time"
)

// ICECandidateReceivedEvent represents a icecandidate_received_event
type ICECandidateReceivedEvent struct {
	Candidate *ICECandidate `json:"candidate,omitempty"`
	ReceivedAt time.Time `json:"received_at,omitempty"`
	RoomId string `json:"room_id,omitempty"`
}
