package models

import (
	"time"
)

// PolicyCacheInvalidatedEvent represents a policy_cache_invalidated_event
type PolicyCacheInvalidatedEvent struct {
	Domain string `json:"domain,omitempty"`
	EventId string `json:"event_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	Reason string `json:"reason,omitempty"`
}
