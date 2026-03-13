package models


// StreamChatResponse represents a stream_chat_response
type StreamChatResponse struct {
	Message *Message `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
