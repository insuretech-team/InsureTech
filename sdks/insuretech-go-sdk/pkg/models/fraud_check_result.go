package models

import (
	"time"
)

// FraudCheckResult represents a fraud_check_result
type FraudCheckResult struct {
	ReviewedAt time.Time `json:"reviewed_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	FraudCheckId string `json:"fraud_check_id,omitempty"`
	ClaimId string `json:"claim_id,omitempty"`
	FraudScore float64 `json:"fraud_score,omitempty"`
	RiskFactors []string `json:"risk_factors,omitempty"`
	Flagged bool `json:"flagged,omitempty"`
	ReviewedBy string `json:"reviewed_by,omitempty"`
}
