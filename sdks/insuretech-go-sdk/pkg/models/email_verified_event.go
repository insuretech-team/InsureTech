package models

import (
	"time"
)

// EmailVerifiedEvent represents a email_verified_event
type EmailVerifiedEvent struct {
	EventId string `json:"event_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	Email string `json:"email,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
