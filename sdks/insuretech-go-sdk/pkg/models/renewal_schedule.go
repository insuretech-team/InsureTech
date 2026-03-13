package models

import (
	"time"
)

// RenewalSchedule represents a renewal_schedule
type RenewalSchedule struct {
	GracePeriodDays int `json:"grace_period_days,omitempty"`
	GracePeriodEnd time.Time `json:"grace_period_end,omitempty"`
	RenewedAt time.Time `json:"renewed_at,omitempty"`
	RenewedPolicyId string `json:"renewed_policy_id,omitempty"`
	Id string `json:"id"`
	PolicyId string `json:"policy_id"`
	RenewalPremium *Money `json:"renewal_premium,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	RenewalDueDate time.Time `json:"renewal_due_date"`
	RenewalType *RenewalType `json:"renewal_type"`
	Status interface{} `json:"status"`
}
