package models


// LossTrendSegment represents a loss_trend_segment
type LossTrendSegment struct {
	ChangePercentage float64 `json:"change_percentage,omitempty"`
	CurrentLossRatio float64 `json:"current_loss_ratio,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
	PreviousLossRatio float64 `json:"previous_loss_ratio,omitempty"`
	SegmentKey string `json:"segment_key,omitempty"`
	SegmentValue string `json:"segment_value,omitempty"`
	TrendDirection string `json:"trend_direction,omitempty"`
}
