package models


// KnowledgeBaseSearchRequest represents a knowledge_base_search_request
type KnowledgeBaseSearchRequest struct {
	Query string `json:"query"`
	Category string `json:"category,omitempty"`
	Language string `json:"language,omitempty"`
}
