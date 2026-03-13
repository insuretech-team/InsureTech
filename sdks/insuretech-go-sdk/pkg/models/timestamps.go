package models

import (
	"time"
)

// Timestamps represents a timestamps
type Timestamps struct {
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
