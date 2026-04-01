package models

import (
	"time"
)

// ComplianceLog represents a compliance_log
type ComplianceLog struct {
	Description string `json:"description"`
	EntityId string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	Evidence string `json:"evidence,omitempty"`
	Id string `json:"id"`
	PerformedBy string `json:"performed_by,omitempty"`
	Regulation string `json:"regulation"`
	Status *ComplianceStatus `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Type *ComplianceType `json:"type"`
}
