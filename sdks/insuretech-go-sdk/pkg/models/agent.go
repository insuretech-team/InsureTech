package models

import (
	"time"
)

// Agent represents a agent
type Agent struct {
	AgentId string `json:"agent_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	Email string `json:"email,omitempty"`
	Status *PartnerAgentStatus `json:"status,omitempty"`
	CommissionRate float64 `json:"commission_rate,omitempty"`
	JoinedAt time.Time `json:"joined_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	PartnerId string `json:"partner_id,omitempty"`
	FullName string `json:"full_name,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	NidNumber string `json:"nid_number,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}
