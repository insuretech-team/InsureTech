package models

import (
	"time"
)

// NotificationTemplate represents a notification_template
type NotificationTemplate struct {
	BodyTemplate string `json:"body_template"`
	Channel *NotificationChannel `json:"channel"`
	CreatedAt time.Time `json:"created_at"`
	IsActive bool `json:"is_active"`
	Language string `json:"language"`
	SubjectTemplate string `json:"subject_template,omitempty"`
	TemplateId string `json:"template_id"`
	TemplateName string `json:"template_name"`
	Type *NotificationType `json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
}
