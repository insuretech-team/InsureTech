package models

import (
	"time"
)

// EmailVerifiedEvent represents a email_verified_event
type EmailVerifiedEvent struct {
	Email string `json:"email,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
