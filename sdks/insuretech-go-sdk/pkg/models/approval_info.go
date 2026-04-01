package models

import (
	"time"
)

// ApprovalInfo represents a approval_info
type ApprovalInfo struct {
	ApprovalId string `json:"approval_id,omitempty"`
	ApprovalLevel int `json:"approval_level,omitempty"`
	ApprovedAt time.Time `json:"approved_at,omitempty"`
	ApprovedBy string `json:"approved_by,omitempty"`
	RejectionReason string `json:"rejection_reason,omitempty"`
	Status *ApprovalStatus `json:"status,omitempty"`
}
