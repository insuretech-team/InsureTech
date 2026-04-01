package models

import (
	"time"
)

// OTPVerifiedEvent represents a otp_verified_event
type OTPVerifiedEvent struct {
	Attempts int `json:"attempts,omitempty"`
	EventId string `json:"event_id,omitempty"`
	OtpId string `json:"otp_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
