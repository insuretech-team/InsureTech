package models

import (
	"time"
)

// AccessGrantedEvent represents a access_granted_event
type AccessGrantedEvent struct {
	EventId string `json:"event_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	Domain string `json:"domain,omitempty"`
	Subject string `json:"subject,omitempty"`
	Object string `json:"object,omitempty"`
	Action string `json:"action,omitempty"`
	MatchedRule string `json:"matched_rule,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
}
