package models

import (
	"time"
)

// AIAnalysis represents a ai_analysis
type AIAnalysis struct {
	Result string `json:"result,omitempty"`
	Recommendations []string `json:"recommendations,omitempty"`
	AnalyzedAt time.Time `json:"analyzed_at,omitempty"`
	AnalysisId string `json:"analysis_id,omitempty"`
	AgentId string `json:"agent_id,omitempty"`
	Type *AnalysisType `json:"type,omitempty"`
	SubjectId string `json:"subject_id,omitempty"`
	ConfidenceScore float64 `json:"confidence_score,omitempty"`
}
