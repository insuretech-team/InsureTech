package models


// ChatResponse represents a chat_response
type ChatResponse struct {
	ConversationEnded bool `json:"conversation_ended,omitempty"`
	ConversationId string `json:"conversation_id,omitempty"`
	SuggestedActions []string `json:"suggested_actions,omitempty"`
}
