package models

import (
	"time"
)

// AccessDeniedEvent represents a access_denied_event
type AccessDeniedEvent struct {
	Reason string `json:"reason,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Object string `json:"object,omitempty"`
	Action string `json:"action,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	UserId string `json:"user_id,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	Domain string `json:"domain,omitempty"`
	Subject string `json:"subject,omitempty"`
}
