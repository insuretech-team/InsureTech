package models

import (
	"time"
)

// PolicyRule represents a policy_rule
type PolicyRule struct {
	Effect interface{} `json:"effect"`
	Description string `json:"description,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	Subject string `json:"subject"`
	Object string `json:"object"`
	Action string `json:"action"`
	Condition string `json:"condition,omitempty"`
	IsActive bool `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	PolicyId string `json:"policy_id"`
	Domain string `json:"domain"`
}
