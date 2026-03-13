package models

import (
	"time"
)

// CommissionCalculatedEvent represents a commission_calculated_event
type CommissionCalculatedEvent struct {
	CommissionId string `json:"commission_id,omitempty"`
	PartnerId string `json:"partner_id,omitempty"`
	AgentId string `json:"agent_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	CommissionAmount *Money `json:"commission_amount,omitempty"`
	CommissionType string `json:"commission_type,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	EventId string `json:"event_id,omitempty"`
}
