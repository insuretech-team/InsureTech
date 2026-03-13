package models


// KnowledgeBaseArticleUpdateRequest represents a knowledge_base_article_update_request
type KnowledgeBaseArticleUpdateRequest struct {
	ArticleId string `json:"article_id"`
	Article *KnowledgeBaseArticle `json:"article,omitempty"`
}
