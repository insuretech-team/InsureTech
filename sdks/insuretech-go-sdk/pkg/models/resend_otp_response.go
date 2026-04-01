package models

import (
	"time"
)

// ResendOTPResponse represents a resend_otp_response
type ResendOTPResponse struct {
	AttemptsRemaining int `json:"attempts_remaining,omitempty"`
	CanRetryAt time.Time `json:"can_retry_at,omitempty"`
	CooldownSeconds int `json:"cooldown_seconds,omitempty"`
	ExpiresInSeconds int `json:"expires_in_seconds,omitempty"`
	OtpId string `json:"otp_id,omitempty"`
	SenderId string `json:"sender_id,omitempty"`
}
