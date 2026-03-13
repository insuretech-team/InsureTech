package models

import (
	"time"
)

// ComplianceLog represents a compliance_log
type ComplianceLog struct {
	EntityId string `json:"entity_id"`
	Description string `json:"description"`
	Evidence string `json:"evidence,omitempty"`
	PerformedBy string `json:"performed_by,omitempty"`
	Type *ComplianceType `json:"type"`
	Regulation string `json:"regulation"`
	Status *ComplianceStatus `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Id string `json:"id"`
	EntityType string `json:"entity_type"`
}
