package models

import (
	"time"
)

// ResendOTPResponse represents a resend_otp_response
type ResendOTPResponse struct {
	Message string `json:"message,omitempty"`
	ExpiresInSeconds int `json:"expires_in_seconds,omitempty"`
	CooldownSeconds int `json:"cooldown_seconds,omitempty"`
	AttemptsRemaining int `json:"attempts_remaining,omitempty"`
	CanRetryAt time.Time `json:"can_retry_at,omitempty"`
	SenderId string `json:"sender_id,omitempty"`
	Error *Error `json:"error,omitempty"`
	OtpId string `json:"otp_id,omitempty"`
}
