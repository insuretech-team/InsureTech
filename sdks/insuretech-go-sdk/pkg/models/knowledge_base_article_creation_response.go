package models


// KnowledgeBaseArticleCreationResponse represents a knowledge_base_article_creation_response
type KnowledgeBaseArticleCreationResponse struct {
	ArticleId string `json:"article_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error *Error `json:"error,omitempty"`
}
