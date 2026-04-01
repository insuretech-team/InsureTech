package models

import (
	"time"
)

// MetricDefinition represents a metric_definition
type MetricDefinition struct {
	CreatedAt time.Time `json:"created_at"`
	Description string `json:"description,omitempty"`
	Dimensions []string `json:"dimensions,omitempty"`
	MetricId string `json:"metric_id"`
	MetricName string `json:"metric_name"`
	Type *MetricType `json:"type"`
	Unit string `json:"unit,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}
