package models

import (
	"time"
)

// DashboardAccessedEvent represents a dashboard_accessed_event
type DashboardAccessedEvent struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	DashboardName string `json:"dashboard_name,omitempty"`
	DashboardType string `json:"dashboard_type,omitempty"`
	EventId string `json:"event_id,omitempty"`
	IpAddress string `json:"ip_address,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
