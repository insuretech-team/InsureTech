package models


// KnowledgeBaseSearchRequest represents a knowledge_base_search_request
type KnowledgeBaseSearchRequest struct {
	Category string `json:"category,omitempty"`
	Language string `json:"language,omitempty"`
	Query string `json:"query"`
}
