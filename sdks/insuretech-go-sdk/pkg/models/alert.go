package models

import (
	"time"
)

// Alert represents a alert
type Alert struct {
	AlertId string `json:"alert_id,omitempty"`
	AlertType *AlertType `json:"alert_type,omitempty"`
	BusinessId string `json:"business_id,omitempty"`
	Count int `json:"count,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	IsRead bool `json:"is_read,omitempty"`
	Message string `json:"message,omitempty"`
	ReadAt time.Time `json:"read_at,omitempty"`
	Severity *AlertSeverity `json:"severity,omitempty"`
	Title string `json:"title,omitempty"`
}
