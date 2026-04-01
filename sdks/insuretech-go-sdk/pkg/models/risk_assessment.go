package models

import (
	"time"
)

// RiskAssessment represents a risk_assessment
type RiskAssessment struct {
	AssessedAt time.Time `json:"assessed_at,omitempty"`
	AssessmentId string `json:"assessment_id,omitempty"`
	DeviceId string `json:"device_id,omitempty"`
	Factors []*RiskFactor `json:"factors,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	Recommendations []string `json:"recommendations,omitempty"`
	RiskLevel *IotRiskLevel `json:"risk_level,omitempty"`
	RiskScore float64 `json:"risk_score,omitempty"`
}
