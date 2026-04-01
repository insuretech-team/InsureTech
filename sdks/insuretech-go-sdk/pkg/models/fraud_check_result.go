package models

import (
	"time"
)

// FraudCheckResult represents a fraud_check_result
type FraudCheckResult struct {
	ClaimId string `json:"claim_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Flagged bool `json:"flagged,omitempty"`
	FraudCheckId string `json:"fraud_check_id,omitempty"`
	FraudScore float64 `json:"fraud_score,omitempty"`
	ReviewedAt time.Time `json:"reviewed_at,omitempty"`
	ReviewedBy string `json:"reviewed_by,omitempty"`
	RiskFactors []string `json:"risk_factors,omitempty"`
}
