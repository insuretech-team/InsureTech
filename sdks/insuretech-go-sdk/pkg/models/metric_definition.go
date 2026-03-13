package models

import (
	"time"
)

// MetricDefinition represents a metric_definition
type MetricDefinition struct {
	Unit string `json:"unit,omitempty"`
	Dimensions []string `json:"dimensions,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	MetricId string `json:"metric_id"`
	MetricName string `json:"metric_name"`
	Description string `json:"description,omitempty"`
	Type *MetricType `json:"type"`
}
