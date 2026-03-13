package models

import (
	"time"
)

// PongResponse represents a pong_response
type PongResponse struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
}
