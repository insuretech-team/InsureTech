package models

import (
	"time"
)

// Commission represents a commission
type Commission struct {
	AgentId string `json:"agent_id,omitempty"`
	CommissionAmount *Money `json:"commission_amount,omitempty"`
	CommissionId string `json:"commission_id,omitempty"`
	CommissionRate float64 `json:"commission_rate,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	PaidAt time.Time `json:"paid_at,omitempty"`
	PartnerId string `json:"partner_id,omitempty"`
	PaymentId string `json:"payment_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	Status *CommissionStatus `json:"status,omitempty"`
	Type *CommissionType `json:"type,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
