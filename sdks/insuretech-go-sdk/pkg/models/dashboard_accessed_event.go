package models

import (
	"time"
)

// DashboardAccessedEvent represents a dashboard_accessed_event
type DashboardAccessedEvent struct {
	EventId string `json:"event_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	DashboardName string `json:"dashboard_name,omitempty"`
	DashboardType string `json:"dashboard_type,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
}
