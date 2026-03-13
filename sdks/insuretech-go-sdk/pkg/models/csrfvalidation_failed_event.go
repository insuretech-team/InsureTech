package models

import (
	"time"
)

// CSRFValidationFailedEvent represents a csrfvalidation_failed_event
type CSRFValidationFailedEvent struct {
	UserAgent string `json:"user_agent,omitempty"`
	RequestMethod string `json:"request_method,omitempty"`
	UserId string `json:"user_id,omitempty"`
	ReceivedCsrfTokenHash string `json:"received_csrf_token_hash,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	RequestPath string `json:"request_path,omitempty"`
	EventId string `json:"event_id,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	ExpectedCsrfTokenHash string `json:"expected_csrf_token_hash,omitempty"`
}
