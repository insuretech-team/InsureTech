package models


// KnowledgeBaseArticleUpdateResponse represents a knowledge_base_article_update_response
type KnowledgeBaseArticleUpdateResponse struct {
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
