package models

import (
	"time"
)

// Dashboard represents a dashboard
type Dashboard struct {
	CreatedAt time.Time `json:"created_at,omitempty"`
	DashboardId string `json:"dashboard_id,omitempty"`
	DashboardName string `json:"dashboard_name,omitempty"`
	Description string `json:"description,omitempty"`
	Layout string `json:"layout,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
