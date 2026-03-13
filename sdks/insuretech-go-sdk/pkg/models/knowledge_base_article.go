package models


// KnowledgeBaseArticle represents a knowledge_base_article
type KnowledgeBaseArticle struct {
	IsPublished bool `json:"is_published,omitempty"`
	HelpfulCount int `json:"helpful_count,omitempty"`
	Id string `json:"id"`
	Title string `json:"title"`
	TitleBn string `json:"title_bn,omitempty"`
	Slug string `json:"slug"`
	Content string `json:"content"`
	Tags []string `json:"tags,omitempty"`
	ViewCount int `json:"view_count,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	Category string `json:"category"`
	ContentBn string `json:"content_bn,omitempty"`
}
