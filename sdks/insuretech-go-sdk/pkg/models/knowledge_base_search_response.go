package models


// KnowledgeBaseSearchResponse represents a knowledge_base_search_response
type KnowledgeBaseSearchResponse struct {
	Articles []*KnowledgeBaseArticle `json:"articles,omitempty"`
	Error *Error `json:"error,omitempty"`
}
