package models

import (
	"time"
)

// AIAgent represents a ai_agent
type AIAgent struct {
	ModelName string `json:"model_name"`
	Status interface{} `json:"status"`
	Type *AgentType `json:"type"`
	Capabilities map[string]interface{} `json:"capabilities,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	LastActiveAt time.Time `json:"last_active_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
	AgentId string `json:"agent_id"`
	AgentName string `json:"agent_name"`
}
