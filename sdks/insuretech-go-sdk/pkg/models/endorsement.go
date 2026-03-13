package models

import (
	"time"
)

// Endorsement represents a endorsement
type Endorsement struct {
	RequestedBy string `json:"requested_by"`
	ApprovedBy string `json:"approved_by,omitempty"`
	EffectiveDate time.Time `json:"effective_date"`
	Id string `json:"id"`
	EndorsementNumber string `json:"endorsement_number"`
	Type *EndorsementType `json:"type"`
	Changes string `json:"changes"`
	PremiumRefundRequired bool `json:"premium_refund_required,omitempty"`
	Status interface{} `json:"status"`
	ApprovedAt time.Time `json:"approved_at,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	PolicyId string `json:"policy_id"`
	Reason string `json:"reason"`
	PremiumAdjustment *Money `json:"premium_adjustment,omitempty"`
}
