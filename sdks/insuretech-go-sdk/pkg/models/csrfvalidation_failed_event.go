package models

import (
	"time"
)

// CSRFValidationFailedEvent represents a csrfvalidation_failed_event
type CSRFValidationFailedEvent struct {
	EventId string `json:"event_id,omitempty"`
	ExpectedCsrfTokenHash string `json:"expected_csrf_token_hash,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	ReceivedCsrfTokenHash string `json:"received_csrf_token_hash,omitempty"`
	RequestMethod string `json:"request_method,omitempty"`
	RequestPath string `json:"request_path,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
