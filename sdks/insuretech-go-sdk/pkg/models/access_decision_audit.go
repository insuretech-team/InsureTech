package models

import (
	"time"
)

// AccessDecisionAudit represents a access_decision_audit
type AccessDecisionAudit struct {
	Action string `json:"action"`
	AuditId string `json:"audit_id"`
	DecidedAt time.Time `json:"decided_at"`
	Decision *PolicyEffect `json:"decision"`
	Domain string `json:"domain"`
	IpAddress string `json:"ip_address,omitempty"`
	MatchedRule string `json:"matched_rule,omitempty"`
	Object string `json:"object"`
	SessionId string `json:"session_id,omitempty"`
	Subject string `json:"subject"`
	UserAgent string `json:"user_agent,omitempty"`
	UserId string `json:"user_id"`
}
