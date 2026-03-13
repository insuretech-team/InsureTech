package models


// ChatResponse represents a chat_response
type ChatResponse struct {
	ConversationEnded bool `json:"conversation_ended,omitempty"`
	Error *Error `json:"error,omitempty"`
	ConversationId string `json:"conversation_id,omitempty"`
	Message string `json:"message,omitempty"`
	SuggestedActions []string `json:"suggested_actions,omitempty"`
}
