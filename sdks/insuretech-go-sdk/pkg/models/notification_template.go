package models

import (
	"time"
)

// NotificationTemplate represents a notification_template
type NotificationTemplate struct {
	TemplateId string `json:"template_id"`
	TemplateName string `json:"template_name"`
	Type *NotificationType `json:"type"`
	Channel *NotificationChannel `json:"channel"`
	SubjectTemplate string `json:"subject_template,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsActive bool `json:"is_active"`
	BodyTemplate string `json:"body_template"`
	Language string `json:"language"`
}
