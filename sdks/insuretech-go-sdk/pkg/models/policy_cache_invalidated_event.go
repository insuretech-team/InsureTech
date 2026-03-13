package models

import (
	"time"
)

// PolicyCacheInvalidatedEvent represents a policy_cache_invalidated_event
type PolicyCacheInvalidatedEvent struct {
	EventId string `json:"event_id,omitempty"`
	Domain string `json:"domain,omitempty"`
	Reason string `json:"reason,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
}
