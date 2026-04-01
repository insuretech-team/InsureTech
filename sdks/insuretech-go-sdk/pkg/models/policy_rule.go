package models

import (
	"time"
)

// PolicyRule represents a policy_rule
type PolicyRule struct {
	Action string `json:"action"`
	Condition string `json:"condition,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string `json:"created_by,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	Description string `json:"description,omitempty"`
	Domain string `json:"domain"`
	Effect interface{} `json:"effect"`
	IsActive bool `json:"is_active"`
	Object string `json:"object"`
	PolicyId string `json:"policy_id"`
	Subject string `json:"subject"`
	UpdatedAt time.Time `json:"updated_at"`
}
