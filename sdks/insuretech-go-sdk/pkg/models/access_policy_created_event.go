package models

import (
	"time"
)

// AccessPolicyCreatedEvent represents a access_policy_created_event
type AccessPolicyCreatedEvent struct {
	EventId string `json:"event_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	PolicyName string `json:"policy_name,omitempty"`
	Resource string `json:"resource,omitempty"`
	Effect string `json:"effect,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
}
