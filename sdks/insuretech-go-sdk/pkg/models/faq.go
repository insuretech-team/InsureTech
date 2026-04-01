package models


// FAQ represents a faq
type FAQ struct {
	Answer string `json:"answer"`
	AnswerBn string `json:"answer_bn,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	Category string `json:"category"`
	DisplayOrder int `json:"display_order,omitempty"`
	Id string `json:"id"`
	IsPublished bool `json:"is_published,omitempty"`
	Question string `json:"question"`
	QuestionBn string `json:"question_bn,omitempty"`
	ViewCount int `json:"view_count,omitempty"`
}
