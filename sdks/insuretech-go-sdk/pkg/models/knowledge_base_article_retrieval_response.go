package models


// KnowledgeBaseArticleRetrievalResponse represents a knowledge_base_article_retrieval_response
type KnowledgeBaseArticleRetrievalResponse struct {
	Article *KnowledgeBaseArticle `json:"article,omitempty"`
	Error *Error `json:"error,omitempty"`
}
