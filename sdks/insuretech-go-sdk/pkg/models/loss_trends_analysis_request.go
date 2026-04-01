package models

import (
	"time"
)

// LossTrendsAnalysisRequest represents a loss_trends_analysis_request
type LossTrendsAnalysisRequest struct {
	AnalysisEnd time.Time `json:"analysis_end,omitempty"`
	AnalysisStart time.Time `json:"analysis_start,omitempty"`
	ComparisonPeriods int `json:"comparison_periods,omitempty"`
	LineOfBusiness string `json:"line_of_business,omitempty"`
	ProductId string `json:"product_id"`
	SegmentationFields []string `json:"segmentation_fields,omitempty"`
}
