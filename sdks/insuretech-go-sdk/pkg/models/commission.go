package models

import (
	"time"
)

// Commission represents a commission
type Commission struct {
	PolicyId string `json:"policy_id,omitempty"`
	CommissionAmount *Money `json:"commission_amount,omitempty"`
	CommissionRate float64 `json:"commission_rate,omitempty"`
	Status *CommissionStatus `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	CommissionId string `json:"commission_id,omitempty"`
	PartnerId string `json:"partner_id,omitempty"`
	AgentId string `json:"agent_id,omitempty"`
	Type *CommissionType `json:"type,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	PaidAt time.Time `json:"paid_at,omitempty"`
}
