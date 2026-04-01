package models

import (
	"time"
)

// AccessDeniedEvent represents a access_denied_event
type AccessDeniedEvent struct {
	Action string `json:"action,omitempty"`
	Domain string `json:"domain,omitempty"`
	EventId string `json:"event_id,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	Object string `json:"object,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	Reason string `json:"reason,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	Subject string `json:"subject,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
