package models

import (
	"time"
)

// SecurityEvent represents a security_event
type SecurityEvent struct {
	Action string `json:"action,omitempty"`
	Authorized bool `json:"authorized,omitempty"`
	EventId string `json:"event_id,omitempty"`
	EventType string `json:"event_type,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	ResourceId string `json:"resource_id,omitempty"`
	ResourceType string `json:"resource_type,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
