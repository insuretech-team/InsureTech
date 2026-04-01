package models

import (
	"time"
)

// Endorsement represents a endorsement
type Endorsement struct {
	ApprovedAt time.Time `json:"approved_at,omitempty"`
	ApprovedBy string `json:"approved_by,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	Changes string `json:"changes"`
	EffectiveDate time.Time `json:"effective_date"`
	EndorsementNumber string `json:"endorsement_number"`
	Id string `json:"id"`
	PolicyId string `json:"policy_id"`
	PremiumAdjustment *Money `json:"premium_adjustment,omitempty"`
	PremiumRefundRequired bool `json:"premium_refund_required,omitempty"`
	Reason string `json:"reason"`
	RequestedBy string `json:"requested_by"`
	Status interface{} `json:"status"`
	Type *EndorsementType `json:"type"`
}
