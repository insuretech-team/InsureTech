package models

import (
	"time"
)

// PolicyRuleDeletedEvent represents a policy_rule_deleted_event
type PolicyRuleDeletedEvent struct {
	EventId string `json:"event_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	DeletedBy string `json:"deleted_by,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
}
