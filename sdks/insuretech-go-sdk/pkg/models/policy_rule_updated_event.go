package models

import (
	"time"
)

// PolicyRuleUpdatedEvent represents a policy_rule_updated_event
type PolicyRuleUpdatedEvent struct {
	EventId string `json:"event_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	UpdatedBy string `json:"updated_by,omitempty"`
	ChangedFields map[string]interface{} `json:"changed_fields,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
}
