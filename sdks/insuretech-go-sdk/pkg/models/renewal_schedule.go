package models

import (
	"time"
)

// RenewalSchedule represents a renewal_schedule
type RenewalSchedule struct {
	AuditInfo interface{} `json:"audit_info"`
	GracePeriodDays int `json:"grace_period_days,omitempty"`
	GracePeriodEnd time.Time `json:"grace_period_end,omitempty"`
	Id string `json:"id"`
	PolicyId string `json:"policy_id"`
	RenewalDueDate time.Time `json:"renewal_due_date"`
	RenewalPremium *Money `json:"renewal_premium,omitempty"`
	RenewalType *RenewalType `json:"renewal_type"`
	RenewedAt time.Time `json:"renewed_at,omitempty"`
	RenewedPolicyId string `json:"renewed_policy_id,omitempty"`
	Status interface{} `json:"status"`
}
