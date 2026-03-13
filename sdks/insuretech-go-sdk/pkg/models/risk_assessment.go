package models

import (
	"time"
)

// RiskAssessment represents a risk_assessment
type RiskAssessment struct {
	AssessmentId string `json:"assessment_id,omitempty"`
	DeviceId string `json:"device_id,omitempty"`
	PolicyId string `json:"policy_id,omitempty"`
	AssessedAt time.Time `json:"assessed_at,omitempty"`
	RiskScore float64 `json:"risk_score,omitempty"`
	RiskLevel *IotRiskLevel `json:"risk_level,omitempty"`
	Factors []*RiskFactor `json:"factors,omitempty"`
	Recommendations []string `json:"recommendations,omitempty"`
}
