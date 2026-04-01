package models

import (
	"time"
)

// Conversation represents a conversation
type Conversation struct {
	AgentId string `json:"agent_id"`
	Context map[string]interface{} `json:"context,omitempty"`
	ConversationId string `json:"conversation_id"`
	EndedAt time.Time `json:"ended_at,omitempty"`
	Messages []*Message `json:"messages,omitempty"`
	StartedAt time.Time `json:"started_at"`
	Status interface{} `json:"status"`
	UserId string `json:"user_id"`
}
