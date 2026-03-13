package models

import (
	"time"
)

// BusinessTargetAchievedEvent represents a business_target_achieved_event
type BusinessTargetAchievedEvent struct {
	TargetName string `json:"target_name,omitempty"`
	TargetType string `json:"target_type,omitempty"`
	TargetValue float64 `json:"target_value,omitempty"`
	AchievedValue float64 `json:"achieved_value,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
}
