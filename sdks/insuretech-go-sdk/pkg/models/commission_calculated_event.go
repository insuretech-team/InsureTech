package models

import (
	"time"
)

// CommissionCalculatedEvent represents a commission_calculated_event
type CommissionCalculatedEvent struct {
	AgentId string `json:"agent_id,omitempty"`
	CommissionAmount *Money `json:"commission_amount,omitempty"`
	CommissionId string `json:"commission_id,omitempty"`
	CommissionType string `json:"commission_type,omitempty"`
	EventId string `json:"event_id,omitempty"`
	PartnerId string `json:"partner_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
