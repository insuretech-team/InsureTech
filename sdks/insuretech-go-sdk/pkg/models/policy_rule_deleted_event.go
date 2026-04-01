package models

import (
	"time"
)

// PolicyRuleDeletedEvent represents a policy_rule_deleted_event
type PolicyRuleDeletedEvent struct {
	DeletedBy string `json:"deleted_by,omitempty"`
	EventId string `json:"event_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
}
