package models

import (
	"time"
)

// RevenueShare represents a revenue_share
type RevenueShare struct {
	AuditInfo interface{} `json:"audit_info"`
	GrossPremium *Money `json:"gross_premium,omitempty"`
	Id string `json:"id"`
	InsurerId string `json:"insurer_id"`
	InsurerShare *Money `json:"insurer_share,omitempty"`
	PlatformShare *Money `json:"platform_share,omitempty"`
	PolicyId string `json:"policy_id"`
	RecordedAt time.Time `json:"recorded_at"`
	RevenueModel string `json:"revenue_model"`
	SplitConfig string `json:"split_config,omitempty"`
}
