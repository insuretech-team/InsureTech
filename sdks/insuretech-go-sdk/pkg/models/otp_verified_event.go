package models

import (
	"time"
)

// OTPVerifiedEvent represents a otp_verified_event
type OTPVerifiedEvent struct {
	EventId string `json:"event_id,omitempty"`
	OtpId string `json:"otp_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Attempts int `json:"attempts,omitempty"`
}
