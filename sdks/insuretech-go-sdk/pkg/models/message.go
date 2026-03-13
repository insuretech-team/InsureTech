package models

import (
	"time"
)

// Message represents a message
type Message struct {
	MessageId string `json:"message_id,omitempty"`
	Role *MessageRole `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}
