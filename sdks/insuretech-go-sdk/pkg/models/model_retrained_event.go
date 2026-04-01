package models

import (
	"time"
)

// ModelRetrainedEvent represents a model_retrained_event
type ModelRetrainedEvent struct {
	AgentId string `json:"agent_id,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	EventId string `json:"event_id,omitempty"`
	ModelName string `json:"model_name,omitempty"`
	ModelVersion string `json:"model_version,omitempty"`
	NewAccuracy float64 `json:"new_accuracy,omitempty"`
	OldAccuracy float64 `json:"old_accuracy,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	TrainingSamples int `json:"training_samples,omitempty"`
}
