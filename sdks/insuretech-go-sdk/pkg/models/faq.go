package models


// FAQ represents a faq
type FAQ struct {
	ViewCount int `json:"view_count,omitempty"`
	AuditInfo interface{} `json:"audit_info"`
	QuestionBn string `json:"question_bn,omitempty"`
	AnswerBn string `json:"answer_bn,omitempty"`
	DisplayOrder int `json:"display_order,omitempty"`
	Id string `json:"id"`
	Category string `json:"category"`
	Question string `json:"question"`
	Answer string `json:"answer"`
	IsPublished bool `json:"is_published,omitempty"`
}
