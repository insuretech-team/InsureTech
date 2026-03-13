package models

import (
	"time"
)

// PolicyRuleCreatedEvent represents a policy_rule_created_event
type PolicyRuleCreatedEvent struct {
	PolicyId string `json:"policy_id,omitempty"`
	Subject string `json:"subject,omitempty"`
	Domain string `json:"domain,omitempty"`
	Action string `json:"action,omitempty"`
	Effect *PolicyEffect `json:"effect,omitempty"`
	Condition string `json:"condition,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Object string `json:"object,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
}
