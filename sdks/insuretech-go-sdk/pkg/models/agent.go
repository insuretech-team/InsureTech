package models

import (
	"time"
)

// Agent represents a agent
type Agent struct {
	AgentId string `json:"agent_id,omitempty"`
	CommissionRate float64 `json:"commission_rate,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	Email string `json:"email,omitempty"`
	FullName string `json:"full_name,omitempty"`
	JoinedAt time.Time `json:"joined_at,omitempty"`
	NidNumber string `json:"nid_number,omitempty"`
	PartnerId string `json:"partner_id,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Status *InsuranceAgentStatus `json:"status,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
