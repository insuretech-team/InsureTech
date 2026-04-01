package models

import (
	"time"
)

// LossTrendsAnalysisResponse represents a loss_trends_analysis_response
type LossTrendsAnalysisResponse struct {
	AnalysisDate time.Time `json:"analysis_date,omitempty"`
	OverallTrend map[string]interface{} `json:"overall_trend,omitempty"`
	Segments []*LossTrendSegment `json:"segments,omitempty"`
	Success bool `json:"success,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}
