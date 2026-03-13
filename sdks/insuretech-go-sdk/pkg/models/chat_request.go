package models


// ChatRequest represents a chat_request
type ChatRequest struct {
	AgentId string `json:"agent_id"`
	Message string `json:"message,omitempty"`
	Context map[string]interface{} `json:"context,omitempty"`
	ConversationId string `json:"conversation_id"`
	UserId string `json:"user_id"`
}
