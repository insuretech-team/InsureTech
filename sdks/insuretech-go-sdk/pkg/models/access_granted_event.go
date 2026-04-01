package models

import (
	"time"
)

// AccessGrantedEvent represents a access_granted_event
type AccessGrantedEvent struct {
	Action string `json:"action,omitempty"`
	Domain string `json:"domain,omitempty"`
	EventId string `json:"event_id,omitempty"`
	MatchedRule string `json:"matched_rule,omitempty"`
	Object string `json:"object,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	Subject string `json:"subject,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
