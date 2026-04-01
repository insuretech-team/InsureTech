package models

import (
	"time"
)

// Message represents a message
type Message struct {
	Content string `json:"content,omitempty"`
	MessageId string `json:"message_id,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Role *MessageRole `json:"role,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
