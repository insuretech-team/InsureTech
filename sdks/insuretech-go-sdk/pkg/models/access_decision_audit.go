package models

import (
	"time"
)

// AccessDecisionAudit represents a access_decision_audit
type AccessDecisionAudit struct {
	DecidedAt time.Time `json:"decided_at"`
	AuditId string `json:"audit_id"`
	UserId string `json:"user_id"`
	Subject string `json:"subject"`
	Object string `json:"object"`
	IpAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	SessionId string `json:"session_id,omitempty"`
	Domain string `json:"domain"`
	Action string `json:"action"`
	Decision *PolicyEffect `json:"decision"`
	MatchedRule string `json:"matched_rule,omitempty"`
}
