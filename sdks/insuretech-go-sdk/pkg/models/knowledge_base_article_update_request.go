package models


// KnowledgeBaseArticleUpdateRequest represents a knowledge_base_article_update_request
type KnowledgeBaseArticleUpdateRequest struct {
	Article *KnowledgeBaseArticle `json:"article,omitempty"`
	ArticleId string `json:"article_id"`
}
