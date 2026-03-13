package models


// KnowledgeBaseArticleDeletionResponse represents a knowledge_base_article_deletion_response
type KnowledgeBaseArticleDeletionResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
