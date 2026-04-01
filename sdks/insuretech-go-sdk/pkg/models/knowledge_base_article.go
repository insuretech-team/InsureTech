package models


// KnowledgeBaseArticle represents a knowledge_base_article
type KnowledgeBaseArticle struct {
	AuditInfo interface{} `json:"audit_info"`
	Category string `json:"category"`
	Content string `json:"content"`
	ContentBn string `json:"content_bn,omitempty"`
	HelpfulCount int `json:"helpful_count,omitempty"`
	Id string `json:"id"`
	IsPublished bool `json:"is_published,omitempty"`
	Slug string `json:"slug"`
	Tags []string `json:"tags,omitempty"`
	Title string `json:"title"`
	TitleBn string `json:"title_bn,omitempty"`
	ViewCount int `json:"view_count,omitempty"`
}
