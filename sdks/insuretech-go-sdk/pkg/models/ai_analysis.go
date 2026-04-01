package models

import (
	"time"
)

// AIAnalysis represents a ai_analysis
type AIAnalysis struct {
	AgentId string `json:"agent_id,omitempty"`
	AnalysisId string `json:"analysis_id,omitempty"`
	AnalyzedAt time.Time `json:"analyzed_at,omitempty"`
	ConfidenceScore float64 `json:"confidence_score,omitempty"`
	Recommendations []string `json:"recommendations,omitempty"`
	Result string `json:"result,omitempty"`
	SubjectId string `json:"subject_id,omitempty"`
	Type *AnalysisType `json:"type,omitempty"`
}
