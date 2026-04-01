package models

import (
	"time"
)

// PolicyRuleCreatedEvent represents a policy_rule_created_event
type PolicyRuleCreatedEvent struct {
	Action string `json:"action,omitempty"`
	Condition string `json:"condition,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	Domain string `json:"domain,omitempty"`
	Effect *PolicyEffect `json:"effect,omitempty"`
	EventId string `json:"event_id,omitempty"`
	Object string `json:"object,omitempty"`
	OccurredAt time.Time `json:"occurred_at,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	Subject string `json:"subject,omitempty"`
}
