package models

import (
	"time"
)

// PingRequest represents a ping_request
type PingRequest struct {
	Timestamp time.Time `json:"timestamp"`
}
