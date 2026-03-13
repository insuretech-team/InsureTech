package models

import (
	"time"
)

// CommissionDetail represents a commission_detail
type CommissionDetail struct {
	CommissionAmount *Money `json:"commission_amount,omitempty"`
	CommissionType string `json:"commission_type,omitempty"`
	EarnedAt time.Time `json:"earned_at,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
}
