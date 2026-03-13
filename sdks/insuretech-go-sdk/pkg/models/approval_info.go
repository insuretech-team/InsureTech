package models

import (
	"time"
)

// ApprovalInfo represents a approval_info
type ApprovalInfo struct {
	Status *ApprovalStatus `json:"status,omitempty"`
	ApprovedBy string `json:"approved_by,omitempty"`
	ApprovedAt time.Time `json:"approved_at,omitempty"`
	RejectionReason string `json:"rejection_reason,omitempty"`
	ApprovalLevel int `json:"approval_level,omitempty"`
	ApprovalId string `json:"approval_id,omitempty"`
}
