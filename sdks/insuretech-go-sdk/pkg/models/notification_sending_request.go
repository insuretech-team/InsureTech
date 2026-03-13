package models


// NotificationSendingRequest represents a notification_sending_request
type NotificationSendingRequest struct {
	RecipientId string `json:"recipient_id"`
	Channel *NotificationChannel `json:"channel,omitempty"`
	Subject string `json:"subject,omitempty"`
	Message string `json:"message,omitempty"`
	TemplateId string `json:"template_id"`
	ScheduleAfterSeconds string `json:"schedule_after_seconds,omitempty"`
	Type *NotificationType `json:"type"`
	TemplateData map[string]interface{} `json:"template_data,omitempty"`
	Priority *NotificationPriority `json:"priority,omitempty"`
}
